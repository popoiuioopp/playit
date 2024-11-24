package main

import (
	"os"

	"playit/consumer"
	"playit/controllers"
	"playit/db"
	"playit/models"

	"github.com/labstack/echo/v4"
)

var twitchClientID = os.Getenv("TwitchClientID")
var twitchClientSecret = os.Getenv("TwitchClientSecret")
var ytAPIKey = os.Getenv("YtAPIKey")
var redirectURI = os.Getenv("RedirectURI")
var twitchToken = os.Getenv("TwitchToken")

var config = models.Config{
	TwitchClientID:     twitchClientID,
	TwitchClientSecret: twitchClientSecret,
	YtAPIKey:           ytAPIKey,
	RedirectURI:        redirectURI,
	TwitchToken:        twitchToken,
}

func main() {
	e := echo.New()

	db.ConnectDB()

	go consumer.StartUpConsumes(config)

	e.Static("/static", "public")

	controllers.RegisterAPIRoutes(e, &config)
	controllers.RegisterViewRoutes(e)
	controllers.RegisterWSRoutes(e, &config)

	e.Logger.Fatal(e.Start(":8080"))
}
