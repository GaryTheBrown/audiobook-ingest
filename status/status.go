package status

type Ingest int

const (
	Detecting Ingest = iota
	MultiChoice
	Manual
	Converting
	Verifying
	Verified
	Failed
	APIDown
)
