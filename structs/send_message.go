package structs

import (
	"audiobook-ingest/enum"
	"audiobook-ingest/status"
)

type SendMessageEnvelope struct {
	Type    enum.SendType `json:"Type"`
	ItemID  string        `json:"ItemID"`
	Payload string        `json:"Payload"`
}

type CreateItem struct {
	AddedAt int64         `json:"AddedAt"`
	Title   string        `json:"Title"`
	Path    string        `json:"Path"`
	Image   string        `json:"Image"`
	Status  status.Ingest `json:"Status"`
}

type AcceptToken struct {
	Accept  bool              `json:"Accept"`
	TokenID string            `json:"TokenID"`
	Name    string            `json:"Name"`
	ImgSrc  string            `json:"ImgSrc"`
	Data    map[string]string `json:"Data"`
}

type AddMetadata struct {
	Fields map[string]string `json:"Fields"`
}

type ChangeStatus struct {
	Status status.Ingest `json:"Status"`
}

type DetectingLogs struct {
	Value string `json:"Value"`
	Icon  string `json:"Icon"`
}

type Choice struct {
	TokenID string            `json:"TokenID"`
	ImgSrc  string            `json:"ImgSrc"`
	Data    map[string]string `json:"Data"`
}
type AddChoices struct {
	Choices []Choice `json:"Choices"`
}
