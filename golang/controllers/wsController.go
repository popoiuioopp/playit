package controllers

import (
	"log"
	"net/http"
	"playit/models"
	"playit/realtime"

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
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v\n", err)
		return err
	}
	defer ws.Close()

	realtime.RegisterClient(ws)
	log.Println("Client connected")

	// Listen for client disconnection
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			log.Println("Client disconnected:", err)
			realtime.UnregisterClient(ws)
			break
		}
	}

	return nil
}
