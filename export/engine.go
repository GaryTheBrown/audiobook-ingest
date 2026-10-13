package export

import (
	"audiobook-ingest/config"
	"audiobook-ingest/data"
	"audiobook-ingest/processor"
	"audiobook-ingest/status"
	"log"
)

type ExportEngine struct{}

func NewEngine() *ExportEngine {
	return &ExportEngine{}
}

func (ee *ExportEngine) StartWorker() {
	log.Printf("Spawning root verification background queue consumer pipelines...")
	for i := 0; i < config.ExportWorkers; i++ {
		go ee.workerLoop()
	}
}

func (ee *ExportEngine) workerLoop() {
	for id := range data.ExportQueue {
		item, exists := data.Store.GetItem(id)
		if !exists {
			continue
		}

		if err := processor.Process.Export(item); err != nil {
			item.SetErrorMessage(err)
			item.SetStatus(status.Failed)
			continue
		}

		data.Store.Remove(id)
	}
}
