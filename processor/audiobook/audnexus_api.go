package audiobook

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func (ab *AudioBook) fetchBookDetailsByID(token string) (map[string]string, error) {
	if token == "" {
		return nil, fmt.Errorf("cannot resolve details for an empty metadata token string")
	}

	baseLinkString := "https://api.audnex.us/books/"
	cleanBaseURL := strings.ReplaceAll(baseLinkString, " ", "")

	fullTargetURL := cleanBaseURL + url.PathEscape(strings.TrimSpace(token))

	req, err := http.NewRequest("GET", fullTargetURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "AudiobookIngestStationPipeline/1.0.0")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("external asset index server returned non-200 status code: %d", resp.StatusCode)
	}

	var res struct {
		Title       string `json:"title"`
		Subtitle    string `json:"subtitle"`
		ReleaseDate string `json:"releaseDate"`
		Description string `json:"summary"`
		Authors     []struct {
			Name string `json:"name"`
		} `json:"authors"`
		Narrators []struct {
			Name string `json:"name"`
		} `json:"narrators"`
		Series struct {
			Name     string `json:"name"`
			Sequence string `json:"sequence"`
		} `json:"series"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	choice := map[string]string{
		"Token":       token,
		"Title":       strings.TrimSpace(res.Title),
		"Year":        strings.TrimSpace(res.ReleaseDate),
		"Description": strings.TrimSpace(res.Description),
		"ImageURL":    cleanBaseURL + url.PathEscape(strings.TrimSpace(token)) + "/image",
		"Author":      "Unknown Author",
		"Narrator":    "Unknown Narrator",
	}

	if res.Subtitle != "" {
		choice["Title"] = fmt.Sprintf("%s: %s", choice["Title"], strings.TrimSpace(res.Subtitle))
	}
	if len(choice["Year"]) >= 4 {
		choice["Year"] = choice["Year"][:4]
	}

	var auths []string
	for _, a := range res.Authors {
		if n := strings.TrimSpace(a.Name); n != "" {
			auths = append(auths, n)
		}
	}
	if len(auths) > 0 {
		choice["Author"] = strings.Join(auths, ", ")
	}

	var narrs []string
	for _, n := range res.Narrators {
		if name := strings.TrimSpace(n.Name); name != "" {
			narrs = append(narrs, name)
		}
	}
	if len(narrs) > 0 {
		choice["Narrator"] = strings.Join(narrs, ", ")
	}

	if res.Series.Name != "" {
		if res.Series.Sequence != "" {
			choice["Series"] = fmt.Sprintf("%s (Book %s)", strings.TrimSpace(res.Series.Name), strings.TrimSpace(res.Series.Sequence))
		} else {
			choice["Series"] = strings.TrimSpace(res.Series.Name)
		}
	}

	return choice, nil
}
