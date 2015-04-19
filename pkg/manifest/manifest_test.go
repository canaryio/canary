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

	m, err := GetManifest(ts.URL, 42)
	if err != nil {
		t.Fatal(err)
	}

	if len(m.Targets) != 1 {
		t.Fatalf("%d targets found, but expected 1", len(m.Targets))
	}

	target := m.Targets[0]
	if target.Name != "canary" {
		t.Fatalf("expected name to be equal to 'canary', go %s", target.URL)
	}

	if target.URL != "http://www.canary.io" {
		t.Fatalf("expected URL to be equal to 'http://www.canary.io', got %s", target.URL)
	}

	if target.Interval != 42 {
		t.Fatal("expected Interval to be equal to 42 when undefined in the manifest json, got %d", target.Interval)
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

	m, err := GetManifest(ts.URL, 42)
	if err != nil {
		t.Fatal(err)
	}

	first_target := m.Targets[0]

	if first_target.Interval != 2 {
		t.Fatalf("expected Interval on the first target to be equal to the manifest json definition of 2, got %d", first_target.Interval)
	}

	second_target := m.Targets[1]

	if second_target.Interval != 4 {
		t.Fatalf("expected Interval on the second target to be equal to the manifest json definition of 4, got %d", second_target.Interval)
	}
}

func TestGetManifestWithTags(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := `{
			"targets": [
				{
					"url": "http://www.canary.io",
					"name": "canary"
				},
				{
					"url": "http://www.github.com",
					"name": "github",
					"tags": [ "tag1", "tag2" ]
				}
			]
		}`

		fmt.Fprintf(w, data)
	}

	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	m, err := GetManifest(ts.URL, 42)
	if err != nil {
		t.Fatal(err)
	}

	first_target := m.Targets[0]

	if len(first_target.Tags) != 0 {
		t.Fatalf("expected Tags on the first target to be empty, got %v", first_target.Tags)
	}

	second_target := m.Targets[1]

	if len(second_target.Tags) != 2 {
		t.Fatalf("expected number of Tags on the second target to be equal to the manifest json definition of 2, got %d", len(second_target.Tags))
	} else {
		if second_target.Tags[0] != "tag1" {
			t.Fatalf("expected first element of Tags on the second target to be equal to the manifest json definition of 'tag1', got %s", second_target.Tags[0])
		}

		if second_target.Tags[1] != "tag2" {
			t.Fatalf("expected first element of Tags on the second target to be equal to the manifest json definition of 'tag2', got %s", second_target.Tags[1])
		}
	}
}

func TestGetManifestWithAttributes(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := `{
			"targets": [
				{
					"url": "http://www.canary.io",
					"name": "canary"
				},
				{
					"url": "http://www.github.com",
					"name": "github",
					"attributes": {
						"foo": "bar",
						"baz": "bap"
					}
				}
			]
		}`

		fmt.Fprintf(w, data)
	}

	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	m, err := GetManifest(ts.URL, 42)
	if err != nil {
		t.Fatal(err)
	}

	first_target := m.Targets[0]

	if first_target.Attributes != nil {
		t.Fatalf("expected Attributes on the first target to be empty, got %v", first_target.Attributes)
	}

	second_target := m.Targets[1]

	if second_target.Attributes == nil {
		t.Fatalf("expected Attributes on the second target to be equal to the manifest json definition, got %v", second_target.Attributes)
	} else {
		foo := second_target.Attributes["foo"]

		if foo != "bar" {
			t.Fatalf("expected 'foo' element of Attributes on the second target to be equal to the manifest json definition of 'bar', got %s", foo)
		}

		baz := second_target.Attributes["baz"]
		if baz != "bap" {
			t.Fatalf("expected 'baz' element of Attributes on the second target to be equal to the manifest json definition of 'bap', got %s", baz)
		}
	}
}

func TestGetManifestWithRequestHeaders(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := `{
			"targets": [
				{
					"url": "http://www.canary.io",
					"name": "canary"
				},
				{
					"url": "http://www.github.com",
					"name": "github",
					"requestHeaders": {
						"X-Foo": "bar"
					}
				}
			]
		}`

		fmt.Fprintf(w, data)
	}

	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	m, err := GetManifest(ts.URL, 42)
	if err != nil {
		t.Fatal(err)
	}

	first_target := m.Targets[0]
	
	if first_target.RequestHeaders != nil {
		t.Fatalf("expected RequestHeaders on the first target to be empty, got %v", first_target.RequestHeaders)
	}

	second_target := m.Targets[1]

	if second_target.RequestHeaders == nil {
		t.Fatalf("expected RequestHeaders on the second target to be equal to the manifest json definition, got %v", second_target.RequestHeaders)
	} else {
		foo := second_target.RequestHeaders["X-Foo"]
		
		if foo != "bar" {
			t.Fatalf("expected 'X-Foo' element of RequestHeaders on the second target to be equal to the manifest json definition of 'bar', got %s", foo)
		}
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

	m, err := GetManifest(ts.URL, 42)
	if err != nil {
		t.Fatal(err)
	}

	// m.StartDelays will the same length as m.Targets
	if len(m.Targets) != len(m.StartDelays) {
		t.Fatalf("expected length of m.StartDelays (%d) to match length of m.Targets (%d)", len(m.Targets), len(m.StartDelays))
	}

	// without calling GenerateRampupDelays, StartDelays should all be zero.
	for index, value := range m.StartDelays {
		if value != 0.0 {
			t.Fatalf("Expected initial start delay to be 0.0, got %d for index %d", value, index)
		}
	}

	// Calling GenerateRampupDelays will update m.StartDelays to an even distribution of
	// millisecond delays for the passed Sampling interval (in seconds).
	m.GenerateRampupDelays(10)

	if m.StartDelays[0] != 0.0 {
		t.Fatalf("The first start delay should be 0.0 even after generation, got %d", m.StartDelays[0])
	}

	if m.StartDelays[1] != 2500.0 {
		t.Fatalf("The second start delay should be 2500.0 ms after generation, got %d", m.StartDelays[1])
	}

	if m.StartDelays[2] != 5000.0 {
		t.Fatalf("The second start delay should be 5000.0 ms after generation, got %d", m.StartDelays[2])
	}

	if m.StartDelays[3] != 7500.0 {
		t.Fatalf("The second start delay should be 7500.0 ms after generation, got %d", m.StartDelays[3])
	}
}

func TestGetManifestWithCapturedHeaders(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := `{
			"targets": [
				{
					"url": "http://www.canary.io",
					"name": "canary"
				},
				{
					"url": "http://www.github.com",
					"name": "github",
					"captureHeaders": [ "Server" ]
				}
			]
		}`

		fmt.Fprintf(w, data)
	}

	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	m, err := GetManifest(ts.URL, 42)
	if err != nil {
		t.Fatal(err)
	}

	first_target := m.Targets[0]
	
	if len(first_target.CaptureHeaders) != 0 {
		t.Fatalf("expected CaptureHeaders on the first target to be empty, got %v", first_target.CaptureHeaders)
	}

	second_target := m.Targets[1]

	if len(second_target.CaptureHeaders) != 1 {
		t.Fatalf("expected Attributes on the second target to be equal to the manifest json definition, got %v", second_target.CaptureHeaders)
	} else {
		h := second_target.CaptureHeaders[0]
		if h != "Server" {
			t.Fatalf("expected first element of CaptureHeaders on the second target to be equal to the manifest json definition of 'Server', got %s", h)
		}
	}
}
