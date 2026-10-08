package websocket

import (
	"audiobook-ingest/data"
	"audiobook-ingest/enum"
	"audiobook-ingest/status"
	"audiobook-ingest/structs"
	"encoding/json"
	"sync"

	"golang.org/x/net/websocket"
)

var (
	clients   = make(map[*websocket.Conn]bool)
	clientMux sync.Mutex

	broadcastQueue = make(chan structs.SendMessageEnvelope, 250)
)

func Setup() {
	data.SendMessage = SendMessage
}

func Start() {
	go startOutboundBroadcastWorker()
}

func AddWebsocket(ws *websocket.Conn) {
	for _, item := range data.Store.Iterator() {
		s := item.Status()
		basePayload := structs.CreateItem{
			AddedAt: item.AddedAt(),
			Title:   item.Name(),
			Path:    item.Path(),
			Image:   item.Image(),
			Status:  s,
		}
		_ = sendNumericalEnvelope(ws, item.ID(), enum.CreateItem, basePayload)

		if item.TokenID() != "" {
			acceptPayload := structs.AcceptToken{
				Accept:  true,
				TokenID: item.TokenID(),
			}
			_ = sendNumericalEnvelope(ws, item.ID(), enum.AcceptToken, acceptPayload)
		}

		currentMeta := item.Metadata()
		if len(currentMeta) > 0 {
			_ = sendNumericalEnvelope(ws, item.ID(), enum.AddMetadata, currentMeta)
		}

		funcChoices := func() {
			choicesPayload := structs.AddChoices{
				Choices: item.Choices(),
			}
			if len(choicesPayload.Choices) > 0 {
				_ = sendNumericalEnvelope(ws, item.ID(), enum.AddChoices, choicesPayload)
			}
		}
		funcDetectingLogs := func() {
			detectingLogs := item.DetectingLogs()
			if len(detectingLogs) > 0 {
				for _, detectingLog := range detectingLogs {
					_ = sendNumericalEnvelope(ws, item.ID(), enum.DetectingAddLog, detectingLog)
				}
			}
		}

		switch s {
		case status.Detecting:
			funcChoices()
			funcDetectingLogs()
		case status.MultiChoice:
			funcChoices()
		case status.Manual:
		case status.Converting:
		case status.Verifying:
		case status.Verified:
		case status.APIDown:
		case status.Failed:
		default:
		}
	}

	clientMux.Lock()
	clients[ws] = true
	clientMux.Unlock()
}

func SendMessage(itemID string, sendType enum.SendType, payloadObj any) {
	jsonPayload, err := json.Marshal(payloadObj)
	if err != nil {
		return
	}

	message := structs.SendMessageEnvelope{
		Type:    sendType,
		ItemID:  itemID,
		Payload: string(jsonPayload),
	}

	select {
	case broadcastQueue <- message:
	default:
	}
}

func sendNumericalEnvelope(ws *websocket.Conn, itemID string, sendType enum.SendType, payloadObj any) error {
	jsonPayload, err := json.Marshal(payloadObj)
	if err != nil {
		return err
	}
	envelope := structs.SendMessageEnvelope{
		Type:    sendType,
		ItemID:  itemID,
		Payload: string(jsonPayload),
	}
	envelopeBytes, err := json.Marshal(envelope)
	if err != nil {
		return err
	}
	return websocket.Message.Send(ws, string(envelopeBytes))
}
