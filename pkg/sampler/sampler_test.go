package sampler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSample(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ok")
	}
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	target := Target{
		URL: ts.URL,
	}

	sampler := New()
	sample, err := sampler.Sample(target)
	if err != nil {
		t.Fatal(err)
	}

	if sample.StatusCode != 200 {
		t.Fatalf("Expected sampleStatus == 200, but got %d\n", sample.StatusCode)
	}
}
