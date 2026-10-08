package folder

import (
	"audiobook-ingest/config"
	"audiobook-ingest/data"
	"audiobook-ingest/processor"

	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

func ScanDirectoryRecursive() {
	var wg sync.WaitGroup

	for _, dir := range config.ImportDirs {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()

			filepath.Walk(path, func(walkPath string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() && walkPath != path {
					_ = processor.Process.Path(walkPath)
				}
				return nil
			})
		}(dir)
	}
	wg.Wait()
}

func StartWatcher() {
	for _, dir := range config.ImportDirs {
		go monitorPath(dir)
	}
}

func monitorPath(targetPath string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
	}
	defer watcher.Close()

	err = filepath.Walk(targetPath, func(walkPath string, info os.FileInfo, err error) error {
		if err == nil && info.IsDir() {
			return watcher.Add(walkPath)
		}
		return nil
	})

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Create == fsnotify.Create {
				fileInfo, err := os.Stat(event.Name)
				if err == nil && fileInfo.IsDir() {
					_ = watcher.Add(event.Name)

					go func(newFolder string) {
						time.Sleep(2 * time.Second)

						if data.Store.ItemExists(newFolder) {
							return
						}

						_ = processor.Process.Path(newFolder)
					}(event.Name)
				}
			}
		case _, ok := <-watcher.Errors:
			if !ok {
				return
			}
		}
	}
}
