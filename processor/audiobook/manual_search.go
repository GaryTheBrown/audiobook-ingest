package audiobook

import (
	"audiobook-ingest/config"
	"audiobook-ingest/data"
	"audiobook-ingest/status"
	"fmt"
	"net/url"
	"strings"
)

func (ab *AudioBook) ExecuteManualSearch(item *data.Item, criteria map[string]string) {
	var titleVal, authorVal string

	for k, v := range criteria {
		switch strings.ToLower(k) {
		case "title":
			titleVal = strings.TrimSpace(v)
		case "author":
			authorVal = strings.TrimSpace(v)
		}
	}

	if titleVal == "" {
		item.SetStatus(status.Manual)
		return
	}

	var searchTerms []string
	searchTerms = append(searchTerms, fmt.Sprintf("title=%s", url.QueryEscape(titleVal)))

	if authorVal != "" {
		searchTerms = append(searchTerms, fmt.Sprintf("author=%s", url.QueryEscape(authorVal)))
	}

	combinedQuery := strings.Join(searchTerms, "&")

	choices, err := ab.queryAudiobookshelfHTTP(combinedQuery)
	if err != nil {
		item.SetStatus(status.Manual)
		return
	}

	item.SetChoices(choices)

	switch len(choices) {
	case 0:
		item.SetStatus(status.Manual)
		return
	case 1:
		matchedBook := choices[0]
		item.SetTokenID(matchedBook.TokenID)
		if !config.DisableAutoIngest {
			item.SetStatus(status.Converting)
			return
		}
	default:
	}
	item.SetStatus(status.MultiChoice)
}
