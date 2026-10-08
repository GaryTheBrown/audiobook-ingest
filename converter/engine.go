package converter

import (
	"audiobook-ingest/config"
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
	// for id := range data.ConversionQueue {
	// }
}
