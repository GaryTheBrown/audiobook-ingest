package processor

import (
	"audiobook-ingest/processor/audiobook"
)

var Process Interface

func init() {
	Process = &audiobook.AudioBook{}
}
