package canary

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/canaryio/canary/pkg/sampler"
)

// Manifest represents configuration data.
type Manifest struct {
	Targets []sampler.Target
}

// GetManifest retreives a manifest from a given URL.
func GetManifest(url string) (manifest Manifest, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &manifest)
	if err != nil {
		return
	}

	return
}
