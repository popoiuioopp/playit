package controllers

import (
	"net/http"

	"playit/messages"
	"playit/music"

	"github.com/labstack/echo/v4"
)

// RegisterAPIRoutes sets up the API routes
func RegisterAPIRoutes(e *echo.Echo) {
	e.GET("/api/healthz", healthCheck)
	e.GET("/api/debug/messages", getDebugMessages)
	e.GET("/api/debug/queue", getDebugQueue)
}

func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "UP"})
}

func getDebugMessages(c echo.Context) error {
	return c.JSON(http.StatusOK, messages.GetMessages())
}

func getDebugQueue(c echo.Context) error {
	return c.JSON(http.StatusOK, music.GetMusicQueue())
}
