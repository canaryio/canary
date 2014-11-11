package canary

type Reporter interface {
	Start()
	Stop()
	Ingest(*Sample) error
}
