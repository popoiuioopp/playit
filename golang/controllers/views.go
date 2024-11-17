package controllers

import (
	"github.com/labstack/echo/v4"
)

// RegisterViewRoutes sets up the HTML views routes
func RegisterViewRoutes(e *echo.Echo) {
	e.GET("/", serveIndex)
}

func serveIndex(c echo.Context) error {
	return c.HTML(200, "")
}
