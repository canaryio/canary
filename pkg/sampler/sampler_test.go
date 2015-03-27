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

func TestSampleWithHeaders(t *testing.T) {
	headerName := "X-Request-Id"
	headerVal  := "abcd-1234"
	
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(headerName, headerVal)
		
		fmt.Fprintf(w, "ok")
	}
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	target := Target{
		URL: ts.URL,
		CaptureHeaders: []string{ headerName },
	}

	sampler := New(10)
	sample, err := sampler.Sample(target)
	if err != nil {
		t.Fatal(err)
	}

	if sample.StatusCode != 200 {
		t.Fatalf("Expected sampleStatus == 200, but got %d\n", sample.StatusCode)
	}
	
	if sample.ResponseHeaders[headerName] != headerVal {
		t.Fatalf("Expected header %s to equal %s but got %s", headerName, headerVal, sample.ResponseHeaders[headerName])
	}
}

func TestSampleWithCanonicalizedHeaderName(t *testing.T) {
	headerName := "x-request-id"
	headerVal  := "abcd-1234"
	
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-Id", headerVal)
		
		fmt.Fprintf(w, "ok")
	}
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	target := Target{
		URL: ts.URL,
		CaptureHeaders: []string{ headerName },
	}

	sampler := New()
	sample, err := sampler.Sample(target)
	if err != nil {
		t.Fatal(err)
	}

	if sample.StatusCode != 200 {
		t.Fatalf("Expected sampleStatus == 200, but got %d\n", sample.StatusCode)
	}
	
	if sample.ResponseHeaders[headerName] != headerVal {
		t.Fatalf("Expected header %s to equal %s but got %s", headerName, headerVal, sample.ResponseHeaders[headerName])
	}
}

func TestSampleWithMissingHeader(t *testing.T) {
	headerName := "X-Request-Id"
	
	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ok")
	}
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	target := Target{
		URL: ts.URL,
		CaptureHeaders: []string{ headerName },
	}

	sampler := New()
	sample, err := sampler.Sample(target)
	if err != nil {
		t.Fatal(err)
	}

	if sample.StatusCode != 200 {
		t.Fatalf("Expected sampleStatus == 200, but got %d\n", sample.StatusCode)
	}
	
	if val, ok := sample.ResponseHeaders[headerName]; ok {
		t.Fatalf("Expected header %s with missing value to be empty but was '%+v'", headerName, val)
	}
}
