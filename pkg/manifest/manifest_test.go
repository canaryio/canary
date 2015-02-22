package manifest

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

func TestGetManifestRampup(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := `{
			"targets": [
				{
					"url": "http://www.canary.io",
					"name": "canary"
				},
				{
					"url": "http://www.github.com",
					"name": "github"
				},
				{
					"url": "http://www.google.com",
					"name": "google"
				},
				{
					"url": "http://www.youtube.com",
					"name": "youtube"
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

	// m.StartDelays will the same length as m.Targets
	if len(m.Targets) != len(m.StartDelays) {
		t.Fatal("expected length of m.StartDelays (%d) to match length of m.Targets (%d)", len(m.Targets), len(m.StartDelays))
	}

	// without calling GenerateRampupDelays, StartDelays should all be zero.
	for index, value := range m.StartDelays {
		if value != 0.0 {
			t.Fatal("Expected initial start delay to be 0.0, got %d for index %d", value, index)
		}
	}

	// Calling GenerateRampupDelays will update m.StartDelays to an even distribution of
	// millisecond delays for the passed Sampling interval (in seconds).
	m.GenerateRampupDelays(10)

	if m.StartDelays[0] != 0.0 {
		t.Fatal("The first start delay should be 0.0 even after generation, got %d", m.StartDelays[0])
	}

	if m.StartDelays[1] != 2500.0 {
		t.Fatal("The second start delay should be 2500.0 ms after generation, got %d", m.StartDelays[1])
	}

	if m.StartDelays[2] != 5000.0 {
		t.Fatal("The second start delay should be 5000.0 ms after generation, got %d", m.StartDelays[2])
	}

	if m.StartDelays[3] != 7500.0 {
		t.Fatal("The second start delay should be 7500.0 ms after generation, got %d", m.StartDelays[3])
	}
}
