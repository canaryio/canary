package manifest

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/canaryio/canary/pkg/sampler"
)

// Manifest represents configuration data.
type Manifest struct {
	Targets     []sampler.Target
	StartDelays []float64
	Hash        string
}

// GenerateRampupDelays generates an even distribution of sensor start delays
// based on the passed number of interval seconds and the number of targets.
func (m *Manifest) GenerateRampupDelays(intervalSeconds int) {
	var intervalMilliseconds = float64(intervalSeconds * 1000)

	var chunkSize = float64(intervalMilliseconds / float64(len(m.Targets)))

	for i := 0.0; i < intervalMilliseconds; i = i + chunkSize {
		m.StartDelays[int((i / chunkSize))] = i
	}
}

func (m *Manifest) setHash() {
	jsonTarget, _ := json.Marshal(m)
	hasher := md5.New()
	hasher.Write(jsonTarget)
	m.Hash = hex.EncodeToString(hasher.Sum(nil))
}

// Get retreives a manifest from a given URL.
func Get(url string, defaultInterval int) (manifest Manifest, err error) {
	var stream io.ReadCloser

	if url[:7] == "file://" {
		stream, err = os.Open(url[7:])
	} else {
		resp, e := http.Get(url)
		err = e
		if err != nil {
			return
		}

		stream = resp.Body
	}

	defer stream.Close()

	body, err := ioutil.ReadAll(stream)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &manifest)
	if err != nil {
		return
	}

	// Store the MD5 hash of the raw manifest
	manifest.setHash()

	// Determine whether to use target.Interval or defaultInterval
	// Targets that lack an interval value in JSON will have their value set to zero. in this case,
	// use defaultInterval
	for ind := range manifest.Targets {
		if manifest.Targets[ind].Interval == 0 {
			manifest.Targets[ind].Interval = defaultInterval
		}
		manifest.Targets[ind].SetHash()
	}

	// Initialize manifest.StartDelays to zeros
	manifest.StartDelays = make([]float64, len(manifest.Targets))
	for i := 0; i < len(manifest.Targets); i++ {
		manifest.StartDelays[i] = 0.0
	}

	return
}
