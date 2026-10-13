package audiobook

import (
	"audiobook-ingest/data"
	"audiobook-ingest/enum"
	"audiobook-ingest/processor/audiobook/apis"
	"audiobook-ingest/status"
	"audiobook-ingest/structs"
	"strings"
)

var (
	SendMessage func(itemID string, sendType enum.SendType, payloadObj any)
)

type AudioBook struct{}

func (ab *AudioBook) FinalExt() string {
	return ".m4b"
}

func (ab *AudioBook) IDLabel() string {
	return "Audible ASIN ID"
}

func (ab *AudioBook) APISearchEnabled() bool {
	return true
}

func (ab *AudioBook) EnvKeyNames() []string {
	return []string{}
}

func (ab *AudioBook) EnvInstructions() []string {
	return []string{}
}

func (ab *AudioBook) GetManualSearchFields() []structs.ManualSearchField {
	return []structs.ManualSearchField{
		{Key: "title", Label: "Book Title", Placeholder: "Enter book title (Required)...", Required: true},
		{Key: "author", Label: "Author / Artist", Placeholder: "Enter author name string...", Required: false},
	}
}

func (ab *AudioBook) SetTokenID(item *data.Item, tokenID string, isUserOverride bool) bool {
	trimmedID := strings.TrimSpace(tokenID)

	//TODO CHECK THE LIST IN CHOICES AND IF ANY MATCH THEN USE THATS DATA SINCE IT@S ALL FROM THE SAME API

	metadata, err := apis.FetchBookDetailsByID(trimmedID)
	if err != nil {
		if isUserOverride {
			SendMessage(item.ID(), enum.AcceptToken, structs.AcceptToken{Accept: false})
		}
		return false
	}

	item.SetFinalMetadata(metadata)
	item.SetTokenID(tokenID)

	if SendMessage != nil {
		acceptPayload := structs.AcceptToken{
			Accept:  true,
			TokenID: trimmedID,
			Data:    item.FinalMetadata().ToMap(),
		}

		SendMessage(item.ID(), enum.AcceptToken, acceptPayload)
	}
	item.SetStatus(status.Converting)
	return true
}
