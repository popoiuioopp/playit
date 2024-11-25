package controllers

import (
	"net/http"
	"playit/views"

	"github.com/labstack/echo/v4"
)

// RegisterViewRoutes sets up the view routes
func RegisterViewRoutes(e *echo.Echo) {
	e.GET("/", HomePageHandler)
	e.GET("/register", RegisterPageHandler)
	e.GET("/login", LoginPageHandler)
	e.GET("/:userName", MusicQueueHandler)
}

func HomePageHandler(c echo.Context) error {
	return Render(c, http.StatusOK, views.HomePage())
}

func RegisterPageHandler(c echo.Context) error {
	return Render(c, http.StatusOK, views.RegisterPage("Register"))
}

func LoginPageHandler(c echo.Context) error {
	return Render(c, http.StatusOK, views.LoginPage())
}

func MusicQueueHandler(c echo.Context) error {
	userName := c.Param("userName")
	return Render(c, http.StatusOK, views.QueuePage("Music Request App", userName))
}
