package websocket

import (
	"audiobook-ingest/config"
	"audiobook-ingest/data"
	"audiobook-ingest/enum"
	"audiobook-ingest/processor"
	"audiobook-ingest/status"
	"audiobook-ingest/structs"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/websocket"
)

func HandleWebsocketStream(ws *websocket.Conn) {
	for {
		var incomingMsg string
		if err := websocket.Message.Receive(ws, &incomingMsg); err != nil {
			clientMux.Lock()
			delete(clients, ws)
			clientMux.Unlock()
			ws.Close()
			break
		}

		var envelope structs.RecievedMessageEnvelope
		if err := json.Unmarshal([]byte(incomingMsg), &envelope); err != nil {
			continue
		}

		handleInboundWebSocketAction(ws, envelope)
	}
}

func handleInboundWebSocketAction(ws *websocket.Conn, env structs.RecievedMessageEnvelope) {
	if strings.ToUpper(env.ItemID) != "SYSTEM" {
		item, exists := data.Store.GetItem(env.ItemID)
		if !exists {
			return
		}

		switch env.Type {
		case enum.SubmitToken:
			var payload structs.SubmitToken
			if err := json.Unmarshal([]byte(env.Payload), &payload); err == nil {
				tokenStr := strings.TrimSpace(payload.TokenID)
				if tokenStr == "" {
					denyResponse := structs.AcceptToken{Accept: false, TokenID: ""}
					SendMessage(item.ID(), enum.AcceptToken, denyResponse)
					return
				}
				item.SetTokenID(tokenStr)
				acceptResponse := structs.AcceptToken{Accept: true, TokenID: tokenStr}
				SendMessage(item.ID(), enum.AcceptToken, acceptResponse)
				item.SetStatus(status.Converting)
			} else {
				denyResponse := structs.AcceptToken{Accept: false, TokenID: ""}
				SendMessage(item.ID(), enum.AcceptToken, denyResponse)
			}
		case enum.FailedDelete:
			purgeItemDiskAssets(item)
			data.Store.Remove(env.ItemID)
		case enum.FailedRetry:
			item.RemoveDetactingLogs()
			item.SetStatus(status.Detecting)
		case enum.SwitchToManual:
			item.SetStatus(status.Manual)
		case enum.SelectChoice:
			var payload structs.SelectChoice
			if err := json.Unmarshal([]byte(env.Payload), &payload); err == nil {
				item.SetTokenID(payload.Choice)
				item.SetStatus(status.Converting)
			}
		case enum.ManualSearch:
			var criteria map[string]string
			if err := json.Unmarshal([]byte(env.Payload), &criteria); err == nil {
				go processor.Process.ExecuteManualSearch(item, criteria)
			} else {
				item.SetStatus(status.Manual)
			}
		case enum.Rescan:
			item.RemoveDetactingLogs()
			item.SetStatus(status.Detecting)
		default:
			panic(fmt.Sprintf("unexpected enum.RecievedType: %#v", env.Type))
		}
	} else {
		switch env.Type {
		case enum.Ping:
			//DO NOTHING HAPPY TO STAY AWAKE
		case enum.DetectingStuck:
			log.Printf("DETECTINGSTUCK")
			statusFunc := func(item *data.Item, targetStatus status.Ingest) {
				if targetStatus != status.Detecting {
					statusPayload := structs.ChangeStatus{
						Status: targetStatus,
					}
					_ = sendNumericalEnvelope(ws, item.ID(), enum.ChangeStatus, statusPayload)
				}
			}

			tokenFunc := func(item *data.Item, tokenID string) {
				acceptPayload := structs.AcceptToken{
					Accept:  true,
					TokenID: tokenID,
					Name:    item.Name(),
					ImgSrc:  item.FoundImgSrc(),
					Data:    item.FoundMetadata(),
				}
				_ = sendNumericalEnvelope(ws, item.ID(), enum.AcceptToken, acceptPayload)
			}

			metaFunc := func(item *data.Item) {
				meta := item.Metadata()
				if len(meta) > 0 {
					_ = sendNumericalEnvelope(ws, item.ID(), enum.AddMetadata, meta)
				}
			}

			choicesFunc := func(item *data.Item) bool {
				choices := item.Choices()
				if len(choices) > 0 {
					choicesPayload := structs.AddChoices{Choices: choices}
					_ = sendNumericalEnvelope(ws, item.ID(), enum.AddChoices, choicesPayload)
					return true
				}
				return false
			}

			var stuckItemIDs []string
			if err := json.Unmarshal([]byte(env.Payload), &stuckItemIDs); err == nil && len(stuckItemIDs) > 0 {
				for _, id := range stuckItemIDs {
					item, exists := data.Store.GetItem(id)
					if !exists {
						continue
					}

					currentStatus := item.Status()
					resolvedToken := item.TokenID()

					if resolvedToken != "" {
						tokenFunc(item, resolvedToken)
						statusFunc(item, currentStatus)
						continue
					}

					foundChoices := choicesFunc(item)
					metaFunc(item)

					if currentStatus == status.MultiChoice && !foundChoices {
						currentStatus = status.Manual
					}

					statusFunc(item, currentStatus)
				}
				return
			}
		}
	}
}

func purgeItemDiskAssets(item *data.Item) {
	if config.NoDelete {
		return
	}

	targetPath := item.FullPathName()
	fileInfo, err := os.Stat(targetPath)
	if err != nil {
		return
	}

	if fileInfo.IsDir() {
		os.RemoveAll(targetPath)
	}

	os.Remove(targetPath)

	parentDir := filepath.Dir(targetPath)
	baseFileNameWithoutExt := strings.TrimSuffix(fileInfo.Name(), filepath.Ext(fileInfo.Name()))

	filepath.Walk(parentDir, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !info.IsDir() && path != targetPath {
			if strings.HasPrefix(info.Name(), baseFileNameWithoutExt) {
				_ = os.Remove(path)
			}
		}
		return nil
	})
}
