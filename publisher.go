package canary

// Publisher is the interface that adds the Publish method.
//
// Pubilsh takes a Target, and Sample, and an error, and is
// expected to deliver that data somewhere.
type Publisher interface {
	Publish(Measurement) error
}
