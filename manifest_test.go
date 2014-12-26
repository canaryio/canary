package canary

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetManifest(t *testing.T) {
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
}
