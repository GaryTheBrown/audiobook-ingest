package websocket

import (
	"encoding/json"
	"log"

	"golang.org/x/net/websocket"
)

func startOutboundBroadcastWorker() {
	for envelope := range broadcastQueue {
		envelopeBytes, err := json.Marshal(envelope)
		if err != nil {
			continue
		}

		rawMessageString := string(envelopeBytes)

		clientMux.Lock()
		for client := range clients {
			log.Printf("Worker SENDING message: %s", rawMessageString)
			if err := websocket.Message.Send(client, rawMessageString); err != nil {
				client.Close()
				delete(clients, client)
			}
		}
		clientMux.Unlock()
	}
}
