package canary

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetManifestWithoutInterval(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := `{
			"targets": [
				{
					"url": "http://www.canary.io",
					"name": "canary"
				}
			]
		}`

		fmt.Fprintf(w, data)
	}

	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	m, err := GetManifest(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	if len(m.Targets) != 1 {
		t.Fatal("%d targets found, but expected 1", len(m.Targets))
	}

	target := m.Targets[0]
	if target.Name != "canary" {
		t.Fatal("expected name to be equal to 'canary', go %s", target.URL)
	}

	if target.URL != "http://www.canary.io" {
		t.Fatal("expected URL to be equal to 'http://www.canary.io', got %s", target.URL)
	}

	if target.Interval != 0 {
		t.Fatal("expected Interval to be equal to zero when undefined in the manifest json, got %d", target.Interval)
	}
}

func TestGetManifestWithInterval(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := `{
			"targets": [
				{
					"url": "http://www.canary.io",
					"name": "canary",
					"interval": 2
				},
				{
					"url": "http://www.github.com",
					"name": "github",
					"interval": 4
				}
			]
		}`

		fmt.Fprintf(w, data)
	}

	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	m, err := GetManifest(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	first_target := m.Targets[0]

	if first_target.Interval != 2 {
		t.Fatal("expected Interval on the first target to be equal to the manifest json definition of 2, got %d", first_target.Interval)
	}

	second_target := m.Targets[1]

	if second_target.Interval != 4 {
		t.Fatal("expected Interval on the second target to be equal to the manifest json definition of 4, got %d", second_target.Interval)
	}
}
