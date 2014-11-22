package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/canaryio/canary"
	"github.com/canaryio/canary/pkg/breadboard"
)

func getManifest(url string) (*canary.Manifest, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http_status=%d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	m := &canary.Manifest{}
	err = json.Unmarshal(body, m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func main() {
	manifestURL := os.Getenv("MANIFEST_URL")
	if manifestURL == "" {
		log.Fatal("MANIFEST_URL not set in ENV")
	}

	m, err := getManifest(manifestURL)
	if err != nil {
		log.Fatal(err)
	}

	// source
	source, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	// breadboard
	b := breadboard.New(source)
	go b.Start()
	b.Update(m)

	select {}
}
