package controllers

import (
	"net/http"
	"playit/views"

	"github.com/labstack/echo/v4"
)

// RegisterViewRoutes sets up the view routes
func RegisterViewRoutes(e *echo.Echo) {
	e.GET("/:userName", MusicQueueHandler)
}

// MusicQueueHandler renders the MusicQueue component.
func MusicQueueHandler(c echo.Context) error {
	userName := c.Param("userName")
	return Render(c, http.StatusOK, views.HomePage("Music Request App", userName))
}
