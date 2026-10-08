package audiobook

import (
	"audiobook-ingest/data"
	"audiobook-ingest/processor/audiobook/structs"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

func (ab *AudioBook) fetchBookDetailsByID(item *data.Item, token string) (map[string]string, error) {
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

	var res structs.AudiobookAudNexusReturn

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	var tagRegexp = regexp.MustCompile(`<[^>]*>`)
	item.SetFinalMetadata(res)
	res.Summary = tagRegexp.ReplaceAllString(res.Summary, "")
	choice := map[string]string{
		"Copyright":        fmt.Sprintf("%d", res.Copyright),
		"Description":      res.Description,
		"FormatType":       res.FormatType,
		"Image":            res.Image,
		"IsAdult":          fmt.Sprintf("%t", res.IsAdult),
		"ISBN":             res.ISBN,
		"Language":         res.Language,
		"LiteratureType":   res.LiteratureType,
		"PublisherName":    res.PublisherName,
		"Rating":           res.Rating,
		"Region":           res.Region,
		"ReleaseDate":      res.ReleaseDate,
		"RuntimeLengthMin": fmt.Sprintf("%d", res.RuntimeLengthMin),
		"Summary":          res.Summary,
		"Title":            res.Title,
	}

	var tmpAuthors []string
	for _, v := range res.Authors {
		tmpAuthors = append(tmpAuthors, v.Name)
	}
	choice["Authors"] = strings.Join(tmpAuthors, ",")
	var tmpGenres []string
	for _, v := range res.Genres {
		tmpGenres = append(tmpGenres, v.Name)
	}
	choice["Genres"] = strings.Join(tmpGenres, ",")
	var tmpNarrators []string
	for _, v := range res.Narrators {
		tmpNarrators = append(tmpNarrators, v.Name)
	}
	choice["Narrators"] = strings.Join(tmpNarrators, ",")

	return choice, nil
}
