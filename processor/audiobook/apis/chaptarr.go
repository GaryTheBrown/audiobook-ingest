package apis

import (
	"audiobook-ingest/config"
	"audiobook-ingest/data"
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func TriggerChaptarrIngest(item *data.Item, filePath string) error {
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
