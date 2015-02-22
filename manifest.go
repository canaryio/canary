package canary


// Manifest is the interface that defines a list of targets
// as well as the startup delays for each of the targets.
//
// GetManifest takes a string location of the manifest to load, and returns a Manifest
// GenerateRampupDelays expects an integer (number of interval seconds) and updates the manifests StartDelays
type Manifest interface {
	GetManifest(string) (Manifest, error)
	GenerateRampupDelays(int)
}
