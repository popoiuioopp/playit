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

var channelSubscriptions = struct {
	sync.RWMutex
	Subscriptions map[string][]string // Key: "platform:channelID", Value: slice of performer usernames
}{
	Subscriptions: make(map[string][]string),
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

func BroadcastMessage(channelID string, platform string, queue []models.SongRequest) {
	// Render the MusicCard component
	var buf bytes.Buffer
	if err := views.MusicCard(queue).Render(context.Background(), &buf); err != nil {
		log.Printf("Error rendering MusicCard component: %v\n", err)
		return
	}
	htmlContent := buf.String()

	// Build the key for channelSubscriptions
	key := platform + ":" + channelID

	// Get the list of performers subscribed to this channel/platform
	channelSubscriptions.RLock()
	performers := channelSubscriptions.Subscriptions[key]
	channelSubscriptions.RUnlock()

	if len(performers) == 0 {
		log.Printf("No performers subscribed to channel %s on platform %s\n", channelID, platform)
		return
	}

	// For each performer, send the message to their connected clients
	performerClients.RLock()
	defer performerClients.RUnlock()

	for _, performer := range performers {
		connections, ok := performerClients.Clients[performer]
		if !ok {
			log.Printf("No connected clients for performer %s\n", performer)
			continue
		}
		for _, conn := range connections {
			err := conn.WriteMessage(websocket.TextMessage, []byte(htmlContent))
			if err != nil {
				log.Printf("WebSocket write error for performer %s: %v\n", performer, err)
				conn.Close()
				// Remove the broken connection
				removeConnection(performer, conn)
			} else {
				log.Printf("HTML content sent to client of performer %s for channel %s on platform %s\n", performer, channelID, platform)
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

func SubscribePerformerToChannel(performer, platform, channelID string) {
	key := platform + ":" + channelID

	channelSubscriptions.Lock()
	defer channelSubscriptions.Unlock()

	// Avoid duplicates
	for _, existingPerformer := range channelSubscriptions.Subscriptions[key] {
		if existingPerformer == performer {
			return // Performer already subscribed
		}
	}

	channelSubscriptions.Subscriptions[key] = append(channelSubscriptions.Subscriptions[key], performer)
}
