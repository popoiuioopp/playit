package main

import (
	"log"
	"net/http"
	"os"

	"playit/consumer"
	"playit/messages"
	"playit/music"

	"github.com/labstack/echo/v4"
)

var twitchClientID = os.Getenv("twitchClientID")
var twitchClientSecret = os.Getenv("twitchClientSecret")
var redirectURI = os.Getenv("redirectURI")
var ytAPIKey = os.Getenv("ytAPIKey")
var twitchToken = os.Getenv("twitchToken")

type Config struct {
	twitchClientID     string
	twitchClientSecret string
	ytAPIKey           string
}

var config = Config{
	twitchClientID:     twitchClientID,
	twitchClientSecret: twitchClientSecret,
	ytAPIKey:           ytAPIKey,
}

func main() {
	e := echo.New()

	music.InitMusicQueue()
	go music.ProcessMusicQueue()
	go consumer.StartYouTubeChatListener("UC3H9YWQl2tNpVOa4AYfJexw", ytAPIKey)
	token, err := consumer.GetTwitchCredential(config.twitchClientID, config.twitchClientSecret)
	if err == nil {
		consumer.ConnectAndConsumeTwitchChat("aiiwan", token)
	} else {
		log.Printf("An Error occurs on twitch consumer %v\n", err)
	}

	e.GET("/healthz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "UP"})
	})

	e.GET("/debug/messages", func(c echo.Context) error {
		return c.JSON(http.StatusOK, messages.GetMessages())
	})

	e.GET("/debug/queue", func(c echo.Context) error {
		return c.JSON(http.StatusOK, music.GetMusicQueue())
	})

	e.Logger.Fatal(e.Start(":8080"))
}
