package data

import (
	"audiobook-ingest/config"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	Store StoreStruct

	DetectingQueue    chan string // Step 1: Scanning file metadata layers
	ConversionQueue   chan string // Step 2: Remuxing loose tracks to M4B
	VerificationQueue chan string // Step 3: Final validation of internal file metadata tags
	ExportQueue       chan string // Step 4: Final move & Chaptarr library API notify

)

func init() {
	Store = NewStore()
	DetectingQueue = make(chan string, 1000)
	ConversionQueue = make(chan string, 1000)
	VerificationQueue = make(chan string, 1000)
	ExportQueue = make(chan string, 1000)
}

func Queue(fullPathName, name string, reqConversion bool) {
	if Store.ItemExists(fullPathName) {
		return
	}

	path := "."
	for _, baseDir := range config.ImportDirs {
		if strings.HasPrefix(fullPathName, baseDir) {
			if rel, err := filepath.Rel(baseDir, fullPathName); err == nil && rel != "." {
				path = filepath.Dir(rel)
				break
			}
		}
	}

	if path == "." {
		path = ""
	}

	item := NewItem(fullPathName, name, path, reqConversion)

	Store.Add(item)
}

func parseWorkerEnv(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return i
}
