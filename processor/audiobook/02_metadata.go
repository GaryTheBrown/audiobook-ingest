package audiobook

import (
	"audiobook-ingest/data"
	"audiobook-ingest/helpers"
	"audiobook-ingest/processor/audiobook/apis"
	"audiobook-ingest/status"
	"fmt"
)

const ModalThreshold = 5

func (ab *AudioBook) DetectMetadata(item *data.Item) {
	foundID, foundMeta := ab.checkMetadataJson(item)
	if foundID {
		return
	}

	if !foundMeta {
		foundID, foundMeta = ab.checkID3Tags(item)
		if foundID {
			return
		}
	}

	if !ab.APISearchEnabled() {
		helpers.SetStatus(item, status.Manual)
		return
	}

	if !foundMeta {
		ab.checkNameParsing(item)
	}

	if ab.RunOnlineMetadataSearch(item) {
		return
	}

	helpers.SetStatus(item, status.Manual)
}

func (ab *AudioBook) RunOnlineMetadataSearch(item *data.Item) bool {
	tags := item.Metadata()

	addDetactingLog(item, "🌐", "Querying Audiobookshelf API provider indexing nodes for match candidate lists...")
	choices, err := apis.AudibleSearch(tags["Title"], tags["Author"], "uk")
	if err != nil {
		addDetactingLog(item, "❌", fmt.Sprintf("Audiobookshelf API request connection error encountered: %v", err))
		return false
	}

	switch len(choices) {
	case 0:
		addDetactingLog(item, "🏁", "Audiobookshelf search completed: Zero alternative match candidates discovered.")
		return false

	case 1:

		return ab.SetTokenID(item, choices[0].TokenID(), false)

	default:
		addDetactingLog(item, "💡", fmt.Sprintf("Audiobookshelf found %d match candidates. Diverting to dashboard triage selection grid.", len(choices)))
		setChoices(item, choices)
		helpers.SetStatus(item, status.MultiChoice)
		return true
	}
}
