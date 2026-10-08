package audiobook

import (
	"audiobook-ingest/data"
	"audiobook-ingest/status"
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
		item.SetStatus(status.Manual)
		return
	}

	if !foundMeta {
		ab.checkNameParsing(item)
	}

	if ab.RunOnlineMetadataSearch(item) {
		return
	}

	item.SetStatus(status.Manual)
}
