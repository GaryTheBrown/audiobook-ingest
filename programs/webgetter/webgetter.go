package webgetter

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func WebGetter(url, ToPath, filename string) (string, error) {
	if strings.TrimSpace(url) == "" {
		return "", fmt.Errorf("cannot download cover from an empty URL string")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "AudiobookIngestStationPipeline/1.0.0")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("artwork server rejected download request with status code: %d", resp.StatusCode)
	}

	ext := ".jpg"
	contentType := strings.ToLower(resp.Header.Get("Content-Type"))
	if strings.Contains(contentType, "image/png") {
		ext = ".png"
	} else if strings.Contains(contentType, "image/jpeg") || strings.Contains(contentType, "image/jpg") {
		ext = ".jpg"
	} else {
		lowerURL := strings.ToLower(url)
		if strings.HasSuffix(lowerURL, ".png") {
			ext = ".png"
		}
	}

	localFilePath := filepath.Join(ToPath, filename+ext)

	out, err := os.Create(localFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create image file container on disk: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to stream cover art binary blocks to disk: %w", err)
	}

	return localFilePath, nil
}
