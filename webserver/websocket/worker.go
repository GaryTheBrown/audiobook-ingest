package websocket

import (
	"encoding/json"

	"golang.org/x/net/websocket"
)

func startOutboundBroadcastWorker() {
	for envelope := range broadcastQueue {
		envelopeBytes, err := json.Marshal(envelope)
		if err != nil {
			continue
		}
		// log.Printf("MSG: %v", envelope)

		clientMux.Lock()
		for client := range clients {
			if err := websocket.Message.Send(client, string(envelopeBytes)); err != nil {
				client.Close()
				delete(clients, client)
			}
		}
		clientMux.Unlock()
	}
}
