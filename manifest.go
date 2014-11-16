package canary

import "encoding/json"

type Site struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

type Service struct {
	Name   string          `json:"name"`
	Config json.RawMessage `json:"config"`
}

type Manifest struct {
	Sites    []*Site    `json:"sites"`
	Services []*Service `json:"services"`
}

func GetManifest(url string) (*Manifest, error) {
	return nil, nil
}
