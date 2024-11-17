package realtime

import (
	"bytes"
	"context"
	"log"
	"playit/models"
	"playit/views"
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
func BroadcastMessage(queue []models.SongRequest) {
	// Render the MusicCard component
	var buf bytes.Buffer
	if err := views.MusicCard(queue).Render(context.Background(), &buf); err != nil {
		log.Printf("Error rendering MusicCard component: %v\n", err)
		return
	}
	htmlContent := buf.String()
	log.Printf("Broadcasting HTML content: %s\n", htmlContent)

	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	// Send the rendered HTML to all connected clients
	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, []byte(htmlContent))
		if err != nil {
			log.Printf("WebSocket write error: %v\n", err)
			client.Close()
			delete(clients, client)
		} else {
			log.Println("HTML content sent to client")
		}
	}
}
