package realtime

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var clientsMutex sync.Mutex

// RegisterClient adds a new WebSocket client to the map
func RegisterClient(client *websocket.Conn) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	clients[client] = true
}

// UnregisterClient removes a WebSocket client from the map
func UnregisterClient(client *websocket.Conn) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	delete(clients, client)
}

// BroadcastMessage sends a message to all connected WebSocket clients
func BroadcastMessage(message interface{}) {
	messageJSON, err := json.Marshal(message)
	log.Printf("Broadcasting message: %s\n", string(messageJSON))
	if err != nil {
		log.Printf("Error marshaling message: %v\n", err)
		return
	}

	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, messageJSON)
		if err != nil {
			log.Printf("WebSocket write error: %v\n", err)
			client.Close()
			delete(clients, client)
		} else {
			log.Println("Message sent to client")
		}
	}
}
