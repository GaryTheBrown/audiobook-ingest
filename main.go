package main

import (
	"audiobook-ingest/folder"
	"audiobook-ingest/metadata"
	"audiobook-ingest/webserver"
	"log"
)

func main() {
	log.Println("Ingestor Initialized: Launching pre-flight system scans...")
	webserver.InitTemplates(HTMLFilesystem)
	metadataDetector := metadata.NewDetector()
	metadataDetector.StartWorker()
	// converterEngine := converter.NewEngine()
	// converterEngine.StartWorker()
	folder.ScanDirectoryRecursive()
	go folder.StartWatcher()
	webserver.Start()
}
