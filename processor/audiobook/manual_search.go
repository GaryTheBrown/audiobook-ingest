package audiobook

import (
	"audiobook-ingest/data"
	"audiobook-ingest/helpers"
	"audiobook-ingest/processor/audiobook/apis"
	"audiobook-ingest/status"
)

func (ab *AudioBook) ExecuteManualSearch(item *data.Item, criteria map[string]string) {

	title, titleFound := criteria["title"]
	if !titleFound {
		helpers.SetStatus(item, status.Manual)
		return
	}

	author, authorFound := criteria["author"]
	if !authorFound {
		author = ""
	}

	choices, err := apis.AudibleSearch(title, author, "uk")
	if err != nil {
		helpers.SetStatus(item, status.Manual)
		return
	}

	switch len(choices) {
	case 0:
		helpers.SetStatus(item, status.Manual)
		return
	case 1:
		ab.SetTokenID(item, choices[0].TokenID(), false)
		item.SetFinalMetadata(choices[0])
		helpers.SetStatus(item, status.Converting)
		return
	default:
		setChoices(item, choices)
	}
	helpers.SetStatus(item, status.MultiChoice)
}
