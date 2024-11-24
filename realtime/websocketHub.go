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

var performerClients = struct {
	sync.RWMutex
	Clients map[string][]*websocket.Conn // Key: performer username
}{
	Clients: make(map[string][]*websocket.Conn),
}

// RegisterClientForPerformer associates a WebSocket client with a performer
func RegisterClientForPerformer(client *websocket.Conn, performer string) {
	performerClients.Lock()
	defer performerClients.Unlock()
	performerClients.Clients[performer] = append(performerClients.Clients[performer], client)
}

// UnregisterClientForPerformer removes a WebSocket client from a performer
func UnregisterClient(performer string, conn *websocket.Conn) {
	removeConnection(performer, conn)
}

func BroadcastMessage(channelID string, queue []models.SongRequest) {
	// Render the MusicCard component
	var buf bytes.Buffer
	if err := views.MusicCard(queue).Render(context.Background(), &buf); err != nil {
		log.Printf("Error rendering MusicCard component: %v\n", err)
		return
	}
	htmlContent := buf.String()

	// Lock the performerClients map
	performerClients.RLock()
	defer performerClients.RUnlock()

	// Find all clients connected to this channel ID
	for performer, connections := range performerClients.Clients {
		log.Printf("Checking connections for performer: %s\n", performer)

		for _, conn := range connections {
			err := conn.WriteMessage(websocket.TextMessage, []byte(htmlContent))
			if err != nil {
				log.Printf("WebSocket write error for performer %s: %v\n", performer, err)
				conn.Close()
				// Remove the broken connection
				removeConnection(performer, conn)
			} else {
				log.Printf("HTML content sent to client of performer %s for channel %s\n", performer, channelID)
			}
		}
	}
}

// Helper function to remove a connection from the performerClients map
func removeConnection(performer string, conn *websocket.Conn) {
	performerClients.Lock()
	defer performerClients.Unlock()

	connections := performerClients.Clients[performer]
	for i, c := range connections {
		if c == conn {
			// Remove the connection from the slice
			performerClients.Clients[performer] = append(connections[:i], connections[i+1:]...)
			break
		}
	}

	// Clean up if no connections remain for the performer
	if len(performerClients.Clients[performer]) == 0 {
		delete(performerClients.Clients, performer)
	}
}
