package main

import (
	"log"
	"os"

	"playit/consumer"
	"playit/controllers"
	"playit/music"

	"github.com/labstack/echo/v4"
)

var twitchClientID = os.Getenv("twitchClientID")
var twitchClientSecret = os.Getenv("twitchClientSecret")
var ytAPIKey = os.Getenv("ytAPIKey")

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
	go consumer.StartYouTubeChatListener("UC3H9YWQl2tNpVOa4AYfJexw", config.ytAPIKey)
	token, err := consumer.GetTwitchCredential(config.twitchClientID, config.twitchClientSecret)
	if err == nil {
		go consumer.ConnectAndConsumeTwitchChat("aiiwan", token)
	} else {
		log.Printf("An Error occurs on twitch consumer %v\n", err)
	}

	controllers.RegisterAPIRoutes(e)
	controllers.RegisterViewRoutes(e)

	e.Logger.Fatal(e.Start(":8080"))
}
