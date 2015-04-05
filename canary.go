package canary

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/canaryio/canary/pkg/manifest"
	"github.com/canaryio/canary/pkg/sampler"
	"github.com/canaryio/canary/pkg/sensor"
)

type Canary struct {
	Config     Config
	Manifest   manifest.Manifest
	Publishers []Publisher
	Sensors    []sensor.Sensor
	OutputChan chan sensor.Measurement
	ReloadChan chan manifest.Manifest
}

// New returns a pointer to a new Publsher.
func New(publishers []Publisher) *Canary {
	return &Canary{
		Publishers: publishers,
		OutputChan: make(chan sensor.Measurement),
	}
}

func (c *Canary) publishMeasurements() {
	// publish each incoming measurement
	for m := range c.OutputChan {
		for _, p := range c.Publishers {
			p.Publish(m)
		}
	}
}

func (c *Canary) SignalHandler() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)
	signal.Notify(signalChan, syscall.SIGHUP)
	for s := range signalChan {
		switch s {
		case syscall.SIGINT:
			for _, sensor := range c.Sensors {
				sensor.Stop()
			}
			os.Exit(0)
		case syscall.SIGHUP:
			manifest, err := manifest.GetManifest(c.Config.ManifestURL, c.Config.DefaultSampleInterval)
			if err != nil {
				log.Fatal(err)
			}
			// Split reload logic into reloader() as to allow other things to trigger a manifest reload
			c.ReloadChan <- manifest
		}
	}
}

func (c *Canary) reloader() {
	if c.ReloadChan == nil {
		c.ReloadChan = make(chan manifest.Manifest)
	}

	for m := range c.ReloadChan {
		stoppingSensors := []sensor.Sensor{}
		for _, sensor := range c.Sensors {
			found := false
			for _, newTarget := range m.Targets {
				if newTarget.Hash == sensor.Target.Hash {
					found = true
				}
			}
			if !found {
				sensor.Stop()
				stoppingSensors = append(stoppingSensors, sensor)
			}

		}
		for _, sensor := range stoppingSensors {
			<-sensor.StopNotifyChan
		}

		c.Manifest = m
		if c.Config.RampupSensors {
			c.Manifest.GenerateRampupDelays(c.Config.DefaultSampleInterval)
		}
		// Start new sensors:
		c.startSensors()
	}
}

func (c *Canary) startSensors() {
	oldSensors := c.Sensors
	c.Sensors = []sensor.Sensor{} // reset the slice

	// spinup a sensor for each target
	for index, target := range c.Manifest.Targets {
		found := false
		for _, oldSensor := range oldSensors {
			if oldSensor.Target.Hash == target.Hash {
				found = true
			}
		}
		if found {
			for _, oldSensor := range oldSensors {
				if oldSensor.Target.Hash == target.Hash {
					c.Sensors = append(c.Sensors, oldSensor)
				}
			}
		} else {
			timeout := target.Interval
			if timeout > c.Config.MaxSampleTimeout {
				timeout = c.Config.MaxSampleTimeout
			}

			sensor := sensor.Sensor{
				Target:         target,
				C:              c.OutputChan,
				Sampler:        sampler.New(timeout),
				StopChan:       make(chan int, 1),
				IsStopped:      false,
				StopNotifyChan: make(chan bool),
				IsOK:           false,
			}
			c.Sensors = append(c.Sensors, sensor)

			go sensor.Start(c.Manifest.StartDelays[index])
		}
	}
}

func (c *Canary) StartAutoReload(interval time.Duration) {
	t := time.NewTicker(interval)
	for {
		<-t.C
		manifest, err := manifest.GetManifest(c.Config.ManifestURL, c.Config.DefaultSampleInterval)
		if err != nil {
			log.Fatal(err)
		}
		if manifest.Hash != c.Manifest.Hash {
			c.ReloadChan <- manifest
		}
	}
}

func (c *Canary) Run() {
	// create and start sensors
	c.startSensors()
	// start a go routine for watching config reloads
	go c.reloader()
	// start a go routine for measurement publishing.
	go c.publishMeasurements()
}
