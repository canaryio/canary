package canary

type Config struct {
	ManifestURL           string
	DefaultSampleInterval int
	RampupSensors         bool
	PublisherList         []string
}
