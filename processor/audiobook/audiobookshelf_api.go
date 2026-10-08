package audiobook

import (
	"audiobook-ingest/data"
	"audiobook-ingest/status"
	"audiobook-ingest/structs"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type AudiobookshelfMatchResponse struct {
	Asin        string `json:"asin"`
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle"`
	Author      string `json:"author"`
	Narrator    string `json:"narrator"`
	SeriesName  string `json:"seriesName"`
	SeriesOrder string `json:"seriesOrder"`
	PublishYear string `json:"publishYear"`
	Cover       string `json:"cover"`
}

func (ab *AudioBook) queryAudiobookshelfHTTP(searchTerm string) ([]structs.Choice, error) {
	targetURL := fmt.Sprintf("%s/api/search/books?provider=audible&%s", AUDIOBOOKSHELF_URL, searchTerm)

	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", AUDIOBOOKSHELF_API_KEY))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "AudiobookIngestStationPipeline/1.0.0")

	client := &http.Client{Timeout: 6 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, fmt.Errorf("authentication failed: invalid or expired AUDIOBOOKSHELF_API_KEY token code")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("local server rejected request with status marker: %d", resp.StatusCode)
	}

	var rawResults []AudiobookshelfMatchResponse
	if err := json.NewDecoder(resp.Body).Decode(&rawResults); err != nil {
		return nil, err
	}

	var processedChoices []structs.Choice
	for _, res := range rawResults {
		coverArtwork := strings.TrimSpace(res.Cover)
		if coverArtwork != "" && !strings.HasPrefix(coverArtwork, "http://") && !strings.HasPrefix(coverArtwork, "https://") {
			coverArtwork = fmt.Sprintf("%s%s", AUDIOBOOKSHELF_URL, coverArtwork)
		}

		choice := structs.Choice{
			TokenID: strings.TrimSpace(res.Asin),
			ImgSrc:  coverArtwork,
			Data: map[string]string{
				"Title":    strings.TrimSpace(res.Title),
				"Author":   strings.TrimSpace(res.Author),
				"Narrator": strings.TrimSpace(res.Narrator),
				"Year":     strings.TrimSpace(res.PublishYear),
			},
		}

		if choice.TokenID == "" {
			continue
		}
		if choice.Data["Author"] == "" {
			choice.Data["Author"] = "Unknown Author"
		}
		if choice.Data["Narrator"] == "" {
			choice.Data["Narrator"] = "Unknown Narrator"
		}

		if res.Subtitle != "" {
			choice.Data["Title"] = fmt.Sprintf("%s: %s", choice.Data["Title"], strings.TrimSpace(res.Subtitle))
		}

		if res.SeriesName != "" {
			if res.SeriesOrder != "" {
				choice.Data["Series"] = fmt.Sprintf("%s (Book %s)", strings.TrimSpace(res.SeriesName), strings.TrimSpace(res.SeriesOrder))
			} else {
				choice.Data["Series"] = strings.TrimSpace(res.SeriesName)
			}
		}

		processedChoices = append(processedChoices, choice)
	}

	return processedChoices, nil
}

func (ab *AudioBook) RunOnlineMetadataSearch(item *data.Item) bool {
	var searchTerms []string
	tags := item.Metadata()

	if tags["Title"] != "" {
		searchTerms = append(searchTerms, fmt.Sprintf("title=%s", url.QueryEscape(tags["Title"])))
		if tags["Author"] != "" {
			searchTerms = append(searchTerms, fmt.Sprintf("author=%s", url.QueryEscape(tags["Author"])))
		}
	} else {
		for i := 1; i <= 30; i++ {
			if segment, exists := tags[fmt.Sprintf("Segment %d", i)]; exists && segment != "" {
				searchTerms = append(searchTerms, segment)
			}
		}
	}

	if len(searchTerms) == 0 {
		item.AddDetactingLog("🔍", "No query tokens available. Running case-insensitive path segment string parsing heuristics...")
		return false
	}
	combinedQuery := strings.Join(searchTerms, "&")
	item.AddDetactingLog("🌐", "Querying Audiobookshelf API provider indexing nodes for match candidate lists...")
	choices, err := ab.queryAudiobookshelfHTTP(combinedQuery)
	if err != nil {
		item.AddDetactingLog("❌", fmt.Sprintf("Audiobookshelf API request connection error encountered: %v", err))
		return false
	}

	switch len(choices) {
	case 0:
		item.AddDetactingLog("🏁", "Audiobookshelf search completed: Zero alternative match candidates discovered.")
		return false

	case 1:
		return item.SetTokenID(choices[0].TokenID)

	default:
		item.AddDetactingLog("💡", fmt.Sprintf("Audiobookshelf found %d match candidates. Diverting to dashboard triage selection grid.", len(choices)))
		item.SetChoices(choices)
		item.SetStatus(status.MultiChoice)
		return true
	}
}
