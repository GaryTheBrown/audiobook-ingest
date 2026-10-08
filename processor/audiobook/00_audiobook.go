package audiobook

import (
	"audiobook-ingest/structs"
	"os"
	"strings"
)

var (
	AUDIOBOOKSHELF_URL     string
	AUDIOBOOKSHELF_API_KEY string

	apiSearchEnabled bool
)

func init() {
	AUDIOBOOKSHELF_URL = strings.TrimSuffix(os.Getenv("AUDIOBOOKSHELF_URL"), "/")
	AUDIOBOOKSHELF_API_KEY = os.Getenv("AUDIOBOOKSHELF_API_KEY")

	apiSearchEnabled = AUDIOBOOKSHELF_URL != "" && AUDIOBOOKSHELF_API_KEY != ""
}

type AudioBook struct{}

func (ab *AudioBook) FinalExt() string {
	return ".m4b"
}

func (ab *AudioBook) IDLabel() string {
	return "Audible ASIN ID"
}

func (ab *AudioBook) APISearchEnabled() bool {
	return apiSearchEnabled
}

func (ab *AudioBook) EnvKeyNames() []string {
	return []string{
		"AUDIOBOOKSHELF_URL",
		"AUDIOBOOKSHELF_API_KEY",
	}
}

func (ab *AudioBook) EnvInstructions() []string {
	return []string{
		"To search manually for an Audible ASIN ID, open audible.co.uk in your browser, search for the book, and copy the 10-digit alphanumeric token starting with 'B' from the product URL path string.",
	}
}

func (ab *AudioBook) GetManualSearchFields() []structs.ManualSearchField {
	return []structs.ManualSearchField{
		{Key: "title", Label: "Book Title", Placeholder: "Enter book title (Required)...", Required: true},
		{Key: "author", Label: "Author / Artist", Placeholder: "Enter author name string...", Required: false},
	}
}

func (ab *AudioBook) NewConversionState() any {
	return &ConversionState{
		totalSeconds: 0,
	}
}

func (ab *AudioBook) CheckTokenID(tokenID string) (exists bool, name string, imgSrc string, metadata map[string]string) {
	exists = false
	name = ""
	imgSrc = ""
	metadata = map[string]string{}
	if allData, err := ab.fetchBookDetailsByID(tokenID); err == nil {
		for key, value := range allData {
			switch key {
			case "Title":
				name = value
			case "image":
				imgSrc = value
			default:
				metadata[key] = value
			}
		}
	}
	return
}
