package controllers

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"playit/models"
	"playit/realtime"
	"playit/repository"
	"playit/views"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins (update this in production)
	},
}

func RegisterWSRoutes(e *echo.Echo, configs *models.Config) {
	e.GET("/ws/queue", HandleWebSocket)
}

// HandleWebSocket handles the WebSocket connection
func HandleWebSocket(c echo.Context) error {
	performer := c.QueryParam("performer")
	if performer == "" {
		return c.String(http.StatusBadRequest, "Performer not specified")
	}

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v\n", err)
		return err
	}
	defer ws.Close()

	// Register the client for the performer
	realtime.RegisterClientForPerformer(ws, performer)
	log.Printf("Client connected for performer: %s\n", performer)

	// Fetch the initial queue from the database for the performer
	initialQueue, err := repository.GetSongRequestsByUserName(performer)
	if err != nil {
		log.Printf("Error fetching initial queue for performer %s: %v\n", performer, err)
		return err
	}

	// Render the MusicCard component with the initial queue
	var buf bytes.Buffer
	if err := views.MusicCard(initialQueue).Render(context.Background(), &buf); err != nil {
		log.Printf("Error rendering MusicCard component for performer %s: %v\n", performer, err)
	} else {
		htmlContent := buf.String()
		// Send the initial queue to the client
		if err := ws.WriteMessage(websocket.TextMessage, []byte(htmlContent)); err != nil {
			log.Printf("WebSocket write error for performer %s: %v\n", performer, err)
			return err
		}
	}

	// Listen for client disconnection
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Client disconnected for performer: %s, error: %v\n", performer, err)
			realtime.UnregisterClient(performer, ws)
			break
		}
	}

	return nil
}
