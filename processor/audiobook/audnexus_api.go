package audiobook

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type Author struct {
	Token string `json:"asin,omitempty"`
	Name  string `json:"name,omitempty"`
}
type Genre struct {
	Token string `json:"asin,omitempty"`
	Name  string `json:"name,omitempty"`
	Type  string `json:"type,omitempty"`
}
type Narrator struct {
	Name string `json:"name,omitempty"`
}

type AudiobookAudNexusReturn struct {
	Token            string     `json:"asin,omitempty"`
	Authors          []Author   `json:"authors,omitempty"`
	Copyright        int        `json:"copyright,omitempty"`
	Description      string     `json:"description,omitempty"`
	FormatType       string     `json:"formatType,omitempty"`
	Genres           []Genre    `json:"genres,omitempty"`
	Image            string     `json:"image,omitempty"`
	IsAdult          bool       `json:"isAdult,omitempty"`
	ISBN             string     `json:"isbn,omitempty"`
	Language         string     `json:"language,omitempty"`
	LiteratureType   string     `json:"literatureType,omitempty"`
	Narrators        []Narrator `json:"narrators,omitempty"`
	PublisherName    string     `json:"publisherName,omitempty"`
	Rating           string     `json:"rating,omitempty"`
	Region           string     `json:"region,omitempty"`
	ReleaseDate      string     `json:"releaseDate,omitempty"`
	RuntimeLengthMin int        `json:"runtimeLengthMin,omitempty"`
	Summary          string     `json:"summary,omitempty"`
	Title            string     `json:"title,omitempty"`
}

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

	var res AudiobookAudNexusReturn

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	var tagRegexp = regexp.MustCompile(`<[^>]*>`)

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
		"Summary":          tagRegexp.ReplaceAllString(res.Summary, ""),
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
