package main

import (
	"audiobook-ingest/converter"
	"audiobook-ingest/folder"
	"audiobook-ingest/metadata"
	"audiobook-ingest/verifier"
	"audiobook-ingest/webserver"
	"log"
)

func main() {
	log.Println("TODO MAKE THIS TITLE PRETTY")
	webserver.InitTemplates(HTMLFilesystem)
	metadataDetector := metadata.NewDetector()
	metadataDetector.StartWorker()
	converterEngine := converter.NewEngine()
	converterEngine.StartWorker()
	verifierEngine := verifier.NewEngine()
	verifierEngine.StartWorker()
	folder.ScanDirectoryRecursive()
	go folder.StartWatcher()
	webserver.Start()
}
