package converter

import (
	"audiobook-ingest/config"
	"audiobook-ingest/data"
	"audiobook-ingest/processor"
	"audiobook-ingest/status"
)

type Engine struct{}

func NewEngine() *Engine {
	return &Engine{}
}

func (e *Engine) StartWorker() {
	for i := 0; i < config.ConversionWorkers; i++ {
		go e.workerLoop()
	}
}

func (e *Engine) workerLoop() {
	for id := range data.ConversionQueue {
		item, exists := data.Store.GetItem(id)
		if !exists {
			continue
		}

		if err := processor.Process.Convert(item); err != nil {
			item.SetErrorMessage(err)
			item.SetStatus(status.Failed)
			continue
		}
	}
}
