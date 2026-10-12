package audiobook

import (
	"audiobook-ingest/data"
	"audiobook-ingest/processor/audiobook/structs"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	base string = "https://api.audnex.us"
)

func booksURL(token string) string {
	return fmt.Sprintf("%s/books/%s?region=uk", base, url.PathEscape(strings.TrimSpace(token)))
}

func chaptersURL(token string) string {
	return fmt.Sprintf("%s/books/%s/chapters?region=uk", base, url.PathEscape(strings.TrimSpace(token)))
}

func apiCall(url string, inStruct any) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "AudiobookIngestStationPipeline/1.0.0")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("external asset index server returned non-200 status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(inStruct); err != nil {
		return err
	}
	return nil
}

func (ab *AudioBook) fetchBookDetailsByID(item *data.Item, token string) error {
	if token == "" {
		return fmt.Errorf("cannot resolve details for an empty metadata token string")
	}

	var audiobook structs.AudnexusApiBookResponse
	if err := apiCall(booksURL(token), &audiobook); err != nil {
		return err
	}

	var chaptersReturn structs.AudnexusApiChapterResponse
	if err := apiCall(chaptersURL(token), &chaptersReturn); err != nil {
		//DO nothing here
	}
	audiobook.Chapters = chaptersReturn.Chapters
	item.SetName(audiobook.Title)
	item.SetFinalCoverImage(audiobook.Image)
	item.SetFinalMetadata(audiobook.ConvertToAudioBookMetadata())
	return nil
}
