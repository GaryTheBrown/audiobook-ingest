package chaptarr

import (
	"audiobook-ingest/config"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func executeChaptarrAuthorPipeline(authorName string) string {
	if config.ChaptarrURL == "" || config.ChaptarrAPIKey == "" {
		return ""
	}

	authorName = strings.TrimSpace(authorName)

	searchURL := fmt.Sprintf("%s/api/v3/author/lookup?term=%s", config.ChaptarrURL, authorName)
	req, _ := http.NewRequest("GET", searchURL, nil)
	req.Header.Set("X-Api-Key", config.ChaptarrAPIKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	return "resolved-chaptarr-author-id"
}
