package controllers

import (
	"net/http"

	"playit/models"

	"github.com/labstack/echo/v4"
)

// RegisterAPIRoutes sets up the API routes
func RegisterAPIRoutes(e *echo.Echo, configs *models.Config) {
	e.GET("/api/healthz", healthCheck)
}

func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "UP"})
}
