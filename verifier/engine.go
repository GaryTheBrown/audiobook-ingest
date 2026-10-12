package verifier

import (
	"audiobook-ingest/config"
	"audiobook-ingest/data"
	"audiobook-ingest/processor"
	"audiobook-ingest/status"
	"log"
)

type VerificationEngine struct{}

func NewEngine() *VerificationEngine {
	return &VerificationEngine{}
}

func (ve *VerificationEngine) StartWorker() {
	log.Printf("Spawning root verification background queue consumer pipelines...")
	for i := 0; i < config.VerifyWorkers; i++ {
		go ve.workerLoop()
	}
}

func (ve *VerificationEngine) workerLoop() {
	for id := range data.VerificationQueue {
		item, exists := data.Store.GetItem(id)
		if !exists {
			continue
		}

		if err := processor.Process.Verify(item); err != nil {
			item.SetErrorMessage(err)
			item.SetStatus(status.Failed)
			continue
		}

		item.SetStatus(status.Verified)
	}
}
