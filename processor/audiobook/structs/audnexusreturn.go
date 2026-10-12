package structs

import (
	"fmt"
	"strings"
	"time"
)

type AudNexusReleaseDate struct {
	time.Time
}

func (t *AudNexusReleaseDate) UnmarshalJSON(b []byte) error {
	formattedTime, err := time.Parse(time.RFC3339, string(b))
	t.Time = formattedTime
	return err
}

type AudnexusApiPerson struct {
	Name string `json:"name"`
}

type AudnexusApiSeries struct {
	Name     string `json:"name"`
	Position string `json:"position,omitempty"`
}

type AudnexusApiGenre struct {
	Name string `json:"name"`
}

type AudnexusApiChapter struct {
	Title         string `json:"title"`
	StartOffsetMs uint64 `json:"startOffsetMs"`
	LengthMs      uint64 `json:"lengthMs"`
}

type AudnexusApiChapterResponse struct {
	Chapters []AudnexusApiChapter `json:"chapters"`
}

type AudnexusApiBookResponse struct {
	Asin          string               `json:"asin"`
	Authors       []AudnexusApiPerson  `json:"authors"`
	Chapters      []AudnexusApiChapter `json:"chapterInfo,omitempty"`
	Description   string               `json:"description,omitempty"`
	Genres        []AudnexusApiGenre   `json:"genres"`
	Image         string               `json:"image,omitempty"`
	Narrators     []AudnexusApiPerson  `json:"narrators"`
	PublisherName string               `json:"publisherName,omitempty"`
	ReleaseDate   AudNexusReleaseDate  `json:"releaseDate,omitempty"`
	Series        []AudnexusApiSeries  `json:"series"`
	Subtitle      string               `json:"subtitle,omitempty"`
	Title         string               `json:"title"`
}

type AudioBookMetadata struct {
	Authors        []string
	Chapters       []AudnexusApiChapter
	Description    string
	Genres         []string
	Narrators      []string
	Publisher      string
	SeriesName     string
	SeriesPosition string
	Subtitle       string
	Year           uint32
}

func (abm *AudioBookMetadata) ExpectedDuration() time.Duration {
	var totalMs uint64 = 0
	for _, ch := range abm.Chapters {
		totalMs += ch.LengthMs
	}
	return time.Duration(totalMs) * time.Millisecond
}

func (r *AudnexusApiBookResponse) ConvertToAudioBookMetadata() *AudioBookMetadata {
	authors := make([]string, len(r.Authors))
	for i, author := range r.Authors {
		authors[i] = author.Name
	}

	narrators := make([]string, len(r.Narrators))
	for i, narrator := range r.Narrators {
		narrators[i] = narrator.Name
	}

	genres := make([]string, len(r.Genres))
	for i, genre := range r.Genres {
		genres[i] = genre.Name
	}

	var seriesName string
	var seriesPosition string
	if len(r.Series) > 0 {
		seriesName = r.Series[0].Name
		if r.Series[0].Position != "" {
			seriesPosition = r.Series[0].Position
		}
	}

	return &AudioBookMetadata{
		Authors:        authors,
		Chapters:       r.Chapters,
		Description:    r.Description,
		Genres:         genres,
		Narrators:      narrators,
		Publisher:      r.PublisherName,
		SeriesName:     seriesName,
		SeriesPosition: seriesPosition,
		Subtitle:       r.Subtitle,
		Year:           uint32(r.ReleaseDate.Year()),
	}
}

func (a AudioBookMetadata) ToMap() map[string]string {
	chapters := make([]string, len(a.Chapters))
	for i, chapter := range a.Chapters {
		chapters[i] = chapter.Title
	}

	return map[string]string{
		"Authors":        strings.Join(a.Authors, ","),
		"Chapters":       strings.Join(chapters, ","),
		"Description":    a.Description,
		"Genres":         strings.Join(a.Genres, ","),
		"Narrators":      strings.Join(a.Narrators, ","),
		"Publisher":      a.Publisher,
		"SeriesName":     a.SeriesName,
		"SeriesPosition": a.SeriesPosition,
		"Subtitle":       a.Subtitle,
		"Year":           fmt.Sprintf("%d", a.Year),
	}
}
