package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

var twitchClientID = os.Getenv("twitchClientID")
var twitchClientSecret = os.Getenv("twitchClientSecret")
var redirectURI = os.Getenv("redirectURI")
var ytAPIKey = os.Getenv("ytAPIKey")

type Config struct {
	Token string
}

var config = Config{}

func storeToken(token string) {
	config.Token = token
}

func main() {
	e := echo.New()

	e.GET("/auth", func(c echo.Context) error {
		authURL := fmt.Sprintf("https://id.twitch.tv/oauth2/authorize?response_type=code&client_id=%s&redirect_uri=%s&scope=chat:read+chat:edit",
			twitchClientID, redirectURI)
		return c.Redirect(http.StatusFound, authURL)
	})

	e.GET("/auth/callback", handleCallback)

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "UP", "token": config.Token})
	})

	e.GET("/debug/messages", func(c echo.Context) error {
		return c.JSON(http.StatusOK, GetMessages())
	})

	e.GET("/debug/queue", func(c echo.Context) error {
		return c.JSON(http.StatusOK, GetMusicQueue())
	})

	go ProcessMusicQueue()
	go StartYouTubeChatListener("UC3H9YWQl2tNpVOa4AYfJexw")

	e.Logger.Fatal(e.Start(":8080"))
}

func handleCallback(c echo.Context) error {
	code := c.QueryParam("code")
	if code == "" {
		return c.String(http.StatusBadRequest, "Code not provided")
	}

	token, err := exchangeCodeForToken(code)
	if err != nil {
		log.Printf("Error exchanging code for token: %v\n", err)
		return c.String(http.StatusInternalServerError, "Error retrieving access token")
	}

	storeToken(token)

	go connectAndConsumeTwitchChat("#midlin_made")

	return c.JSON(http.StatusOK, map[string]string{"access_token": token})
}
