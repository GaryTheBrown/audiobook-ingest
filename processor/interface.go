package processor

import (
	"audiobook-ingest/data"
	"audiobook-ingest/structs"
)

type Interface interface {
	IDLabel() string
	FinalExt() string
	APISearchEnabled() bool
	EnvKeyNames() []string
	EnvInstructions() []string
	GetManualSearchFields() []structs.ManualSearchField
	Path(path string) error
	DetectMetadata(item *data.Item)
	ExecuteManualSearch(item *data.Item, criteria map[string]string)
	CheckTokenID(item *data.Item, tokenID string) error
	Convert(item *data.Item) error
}
