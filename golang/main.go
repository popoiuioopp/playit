package main

import (
	"os"

	"playit/consumer"
	"playit/controllers"
	"playit/models"
	"playit/music"

	"github.com/labstack/echo/v4"
)

var twitchClientID = os.Getenv("twitchClientID")
var twitchClientSecret = os.Getenv("twitchClientSecret")
var ytAPIKey = os.Getenv("ytAPIKey")
var redirectURI = os.Getenv("redirectURI")
var twitchChannelName = "ka_beeja"
var youtubeChannelId = "UC3H9YWQl2tNpVOa4AYfJexw"

var config = models.Config{
	TwitchClientID:     twitchClientID,
	TwitchClientSecret: twitchClientSecret,
	YtAPIKey:           ytAPIKey,
	RedirectURI:        redirectURI,
	TwitchChannelName:  twitchChannelName,
	YoutubeChannelId:   youtubeChannelId,
}

func main() {
	e := echo.New()

	music.InitMusicQueue()
	go music.ProcessMusicQueue()
	go consumer.StartYouTubeChatListener("UC3H9YWQl2tNpVOa4AYfJexw", config.YtAPIKey)

	e.Static("/static", "public")

	controllers.RegisterAPIRoutes(e, &config)
	controllers.RegisterViewRoutes(e)
	controllers.RegisterWSRoutes(e, &config)

	e.Logger.Fatal(e.Start(":8080"))
}
