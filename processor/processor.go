package processor

import (
	"audiobook-ingest/data"
	"audiobook-ingest/processor/audiobook"
)

var Process Interface

func init() {
	Process = &audiobook.AudioBook{}
	data.CheckTokenID = Process.CheckTokenID
}
