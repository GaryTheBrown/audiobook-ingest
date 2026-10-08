package metadata

import (
	"audiobook-ingest/config"
	"audiobook-ingest/data"
	"audiobook-ingest/processor"
	"audiobook-ingest/status"
	"path/filepath"
)

type Detector struct{}

func NewDetector() *Detector {
	return &Detector{}
}

func (d *Detector) StartWorker() {
	for i := 0; i < config.DetectingWorkers; i++ {
		go d.workerLoop(i)
	}
}

func (d *Detector) workerLoop(workerID int) {
	for itemID := range data.DetectingQueue {
		d.ProcessQueueItem(itemID)
	}
}

// 💡⚠✔❌🏁
func (d *Detector) ProcessQueueItem(itemID string) {
	item, exists := data.Store.GetItem(itemID)

	if !exists {
		return
	}
	targetExt := processor.Process.FinalExt()
	if item.RequiredConversion() && filepath.Ext(item.FullPathName()) == targetExt {
		item.AddDetactingLog("⚠", "[SAFETY CATCH] Correcting conversion flag mismatch")
		item.SetRequiredConversion(false)
	}

	if !processor.Process.APISearchEnabled() {
		processor.Process.DetectMetadata(item)
		foundOfflineID := item.TokenID()
		if foundOfflineID == "" {
			item.SetStatus(status.Manual)
			return
		}
	} else {
		processor.Process.DetectMetadata(item)
	}

	resolvedID := item.TokenID()
	currentStatus := item.Status()

	if resolvedID != "" {
		item.SetStatus(status.Converting)
	} else {
		if currentStatus == status.Detecting {
		}
	}
}
