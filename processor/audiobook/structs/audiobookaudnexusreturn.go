package structs

type AudiobookAudNexusReturnAuthor struct {
	Token string `json:"asin,omitempty"`
	Name  string `json:"name,omitempty"`
}
type AudiobookAudNexusReturnGenre struct {
	Token string `json:"asin,omitempty"`
	Name  string `json:"name,omitempty"`
	Type  string `json:"type,omitempty"`
}
type AudiobookAudNexusReturnNarrator struct {
	Name string `json:"name,omitempty"`
}

type AudiobookAudNexusReturn struct {
	Token            string                            `json:"asin,omitempty"`
	Authors          []AudiobookAudNexusReturnAuthor   `json:"authors,omitempty"`
	Copyright        int                               `json:"copyright,omitempty"`
	Description      string                            `json:"description,omitempty"`
	FormatType       string                            `json:"formatType,omitempty"`
	Genres           []AudiobookAudNexusReturnGenre    `json:"genres,omitempty"`
	Image            string                            `json:"image,omitempty"`
	IsAdult          bool                              `json:"isAdult,omitempty"`
	ISBN             string                            `json:"isbn,omitempty"`
	Language         string                            `json:"language,omitempty"`
	LiteratureType   string                            `json:"literatureType,omitempty"`
	Narrators        []AudiobookAudNexusReturnNarrator `json:"narrators,omitempty"`
	PublisherName    string                            `json:"publisherName,omitempty"`
	Rating           string                            `json:"rating,omitempty"`
	Region           string                            `json:"region,omitempty"`
	ReleaseDate      string                            `json:"releaseDate,omitempty"`
	RuntimeLengthMin int                               `json:"runtimeLengthMin,omitempty"`
	Summary          string                            `json:"summary,omitempty"`
	Title            string                            `json:"title,omitempty"`
}
