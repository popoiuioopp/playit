package controllers

import (
	"net/http"
	"playit/views"

	"github.com/labstack/echo/v4"
)

// RegisterViewRoutes sets up the view routes
func RegisterViewRoutes(e *echo.Echo) {
	e.GET("/", HomeHandler)
}

// HomeHandler renders the HomePage component.
func HomeHandler(c echo.Context) error {
	return Render(c, http.StatusOK, views.HomePage("Music Request App"))
}
