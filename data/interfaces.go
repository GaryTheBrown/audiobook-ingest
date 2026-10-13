package data

import "audiobook-ingest/structs"

type MetadataInterface interface {
	ToMap() map[string]string
	ToChoice() structs.Choice
	TokenID() string
}
