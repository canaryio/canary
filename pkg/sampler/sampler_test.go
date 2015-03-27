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

	sampler := New(10)
	sample, err := sampler.Sample(target)
	if err != nil {
		t.Fatal(err)
	}

	if sample.StatusCode != 200 {
		t.Fatalf("Expected sampleStatus == 200, but got %d\n", sample.StatusCode)
	}
}

func TestSampleWithRequestHeaders(t *testing.T) {
	var header http.Header
	
	handler := func(w http.ResponseWriter, r *http.Request) {
		header = r.Header
		
		fmt.Fprintf(w, "ok")
	}
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	target := Target{
		URL: ts.URL,
		RequestHeaders: map[string]string{
			"X-Foo": "bar",
		},
	}

	sampler := New(10)
	sample, err := sampler.Sample(target)
	if err != nil {
		t.Fatal(err)
	}

	if sample.StatusCode != 200 {
		t.Fatalf("Expected sampleStatus == 200, but got %d\n", sample.StatusCode)
	}
	
	h := header.Get("X-Foo")
	if h != "bar" {
		t.Fatalf("Expected request header X-Foo to be 'bar' but was '%s'", h)
	}
}
