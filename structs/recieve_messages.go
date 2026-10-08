package structs

import "audiobook-ingest/enum"

type RecievedMessageEnvelope struct {
	Type    enum.RecievedType `json:"Type"`
	ItemID  string            `json:"ItemID"`
	Payload string            `json:"Payload"`
}

type SubmitToken struct {
	TokenID string `json:"TokenID"`
}

type SelectChoice struct {
	Choice string `json:"Choice"`
}
