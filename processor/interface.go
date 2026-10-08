package processor

import (
	"audiobook-ingest/data"
	"audiobook-ingest/structs"
	"os/exec"
)

type Interface interface {
	IDLabel() string
	FinalExt() string
	APISearchEnabled() bool
	EnvKeyNames() []string
	EnvInstructions() []string
	GetManualSearchFields() []structs.ManualSearchField
	NewConversionState() any
	CompileConversionCmd(item *data.Item) (*exec.Cmd, error)
	ParseProgressLine(line string, state any) (pct int, metaText string, showInLog bool)

	Path(path string) error
	DetectMetadata(item *data.Item)
	ExecuteManualSearch(item *data.Item, criteria map[string]string)
	CheckTokenID(tokenID string) (exists bool, name string, imgSrc string, metadata map[string]string)
}
