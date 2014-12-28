package canary

import "time"

func ExampleStdoutPublisher_Publish() {
	target := Target{
		URL: "http://www.canary.io",
	}

	t1, _ := time.Parse(time.RFC3339, "2014-12-28T00:00:00Z")
	t2, _ := time.Parse(time.RFC3339, "2014-12-28T00:00:01Z")

	sample := Sample{
		T1:         t1,
		T2:         t2,
		StatusCode: 200,
	}

	p := StdoutPublisher{}
	p.Publish(target, sample, nil)
	// Output:
	// 2014-12-28T00:00:01Z http://www.canary.io 200 1000 true
}
