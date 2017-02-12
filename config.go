package canary

import "time"

type Config struct {
	ManifestURL           string
	DefaultSampleInterval int
	RampupSensors         bool
	ReloadInterval        time.Duration
	MaxReloadFailures     int
	MaxSampleTimeout      int
}
