package audiobook

import (
	"audiobook-ingest/config"
	"audiobook-ingest/data"
	"audiobook-ingest/status"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (ab *AudioBook) HandleVerifiedStage(item *data.Item) error {
	bookFolder := item.FullPathName()
	m4bName := fmt.Sprintf("%s.m4b", item.Name())
	srcPath := filepath.Join(bookFolder, m4bName)

	destPath := filepath.Join(config.ExportDir, m4bName)

	if err := moveFileCrossDevice(srcPath, destPath); err != nil {
		item.SetStatus(status.Failed)
		return fmt.Errorf("failed to move finalized file to export directory: %w", err)
	}

	go ab.notifyDownstreamLibraries(item, destPath)

	//TODO HERE NEED TO CLEAN UP FILES IN IMPORT FOLDER
	return nil
}

func (ab *AudioBook) notifyDownstreamLibraries(item *data.Item, finalExportPath string) {
	if config.ChaptarrAPIKey != "" {
		ab.triggerChaptarrIngest(item, finalExportPath)
	}
	// ab.TriggerAudiobookshelfScan()
}

func (ab *AudioBook) triggerChaptarrIngest(item *data.Item, filePath string) error {
	targetURL := fmt.Sprintf("%s/api/v1/ingest", config.ChaptarrURL)

	jsonPayload := fmt.Sprintf(`{"asin":"%s","filePath":"%s","title":"%s"}`, item.TokenID(), strings.ReplaceAll(filePath, "\\", "\\\\"), strings.ReplaceAll(item.Name(), `"`, `\"`))

	req, err := http.NewRequest("POST", targetURL, bytes.NewBuffer([]byte(jsonPayload)))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.ChaptarrAPIKey))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "AudiobookIngestStation/1.0.0")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("chaptarr server rejected notification payload with status: %d", resp.StatusCode)
	}

	return nil
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
