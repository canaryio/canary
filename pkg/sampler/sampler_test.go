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

	sample, err := Ping(target, 1)
	if err != nil {
		t.Fatal(err)
	}

	if sample.StatusCode != 200 {
		t.Fatalf("Expected sampleStatus == 200, but got %d\n", sample.StatusCode)
	}
}

func TestSampleWithHeaders(t *testing.T) {
	headerName := "X-Request-Id"
	headerVal := "abcd-1234"

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(headerName, headerVal)

		fmt.Fprintf(w, "ok")
	}
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	target := Target{
		URL: ts.URL,
	}

	sample, err := Ping(target, 1)
	if err != nil {
		t.Fatal(err)
	}

	if sample.StatusCode != 200 {
		t.Fatalf("Expected sampleStatus == 200, but got %d\n", sample.StatusCode)
	}

	if sample.ResponseHeaders.Get(headerName) != headerVal {
		t.Fatalf("Expected header %s to equal %s but got %s", headerName, headerVal, sample.ResponseHeaders.Get(headerName))
	}
}

func TestSampleWithCanonicalizedHeaderName(t *testing.T) {
	headerName := "x-request-id"
	headerVal := "abcd-1234"

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-Id", headerVal)

		fmt.Fprintf(w, "ok")
	}
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	target := Target{
		URL: ts.URL,
	}

	sample, err := Ping(target, 1)
	if err != nil {
		t.Fatal(err)
	}

	if sample.StatusCode != 200 {
		t.Fatalf("Expected sampleStatus == 200, but got %d\n", sample.StatusCode)
	}

	if sample.ResponseHeaders.Get(headerName) != headerVal {
		t.Fatalf("Expected header %s to equal %s but got %s", headerName, headerVal, sample.ResponseHeaders.Get(headerName))
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
	}

	sample, err := Ping(target, 1)
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

	sample, err := Ping(target, 1)
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
