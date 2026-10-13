package audiobook

import (
	"audiobook-ingest/config"
	"audiobook-ingest/data"
	"audiobook-ingest/processor/audiobook/apis"
	"audiobook-ingest/status"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func (ab *AudioBook) Export(item *data.Item) error {
	bookFolder := item.FullPathName()
	m4bName := fmt.Sprintf("%s.m4b", item.Name())
	srcPath := filepath.Join(bookFolder, m4bName)

	destPath := filepath.Join(config.ExportDir, m4bName)

	if err := moveFileCrossDevice(srcPath, destPath); err != nil {
		item.SetStatus(status.Failed)
		return fmt.Errorf("failed to move finalized file to export directory: %w", err)
	}

	go ab.notifyDownstreamLibraries(item, destPath)

	dirEntries, err := os.ReadDir(bookFolder)
	if err != nil {
	} else {
		var actualRemainingCount int
		for _, entry := range dirEntries {
			if !strings.HasPrefix(entry.Name(), ".") {
				actualRemainingCount++
			}
		}

		if actualRemainingCount == 0 || (actualRemainingCount == 1 && (dirEntries[0].Name() == "cover.jpg" || dirEntries[0].Name() == "cover.png" || dirEntries[0].Name() == "ffmpeg_concat_queue.txt")) {
			if purgeErr := os.RemoveAll(bookFolder); purgeErr != nil {
			}
		} else {
			_ = os.Remove(filepath.Join(bookFolder, "cover.jpg"))
			_ = os.Remove(filepath.Join(bookFolder, "cover.png"))
			_ = os.Remove(filepath.Join(bookFolder, "ffmpeg_concat_queue.txt"))
			_ = os.Remove(filepath.Join(bookFolder, "ffmpeg_metadata_tags.txt"))

			err = filepath.Walk(bookFolder, func(path string, info os.FileInfo, walkErr error) error {
				if walkErr != nil {
					return walkErr
				}

				if info.IsDir() && path != bookFolder {
					return filepath.SkipDir
				}

				if !info.IsDir() {
					ext := strings.ToLower(filepath.Ext(path))
					if ext == ".mp3" || ext == ".m4a" {
						_ = os.Remove(path)
					}
				}
				return nil
			})
			if err != nil {
			}
		}
	}

	return nil
}

func (ab *AudioBook) notifyDownstreamLibraries(item *data.Item, finalExportPath string) {
	if config.ChaptarrAPIKey != "" {
		apis.TriggerChaptarrIngest(item, finalExportPath)
	}
}

func moveFileCrossDevice(src, dest string) error {
	if err := os.Rename(src, dest); err == nil {
		return nil
	}

	input, err := os.Open(src)
	if err != nil {
		return err
	}
	defer input.Close()

	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return err
	}

	output, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer output.Close()

	if _, err = io.Copy(output, input); err != nil {
		return err
	}

	return os.Remove(src)
}
