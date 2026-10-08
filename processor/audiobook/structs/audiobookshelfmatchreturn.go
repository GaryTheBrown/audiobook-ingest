package structs

type AudiobookshelfMatchReturn struct {
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
