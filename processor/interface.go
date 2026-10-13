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
	SetTokenID(item *data.Item, tokenID string, isUserOverride bool) bool
	Convert(item *data.Item) error
	Verify(item *data.Item) error
	HandleVerifiedStage(item *data.Item) error
}
