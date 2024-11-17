package controllers

import (
	"fmt"
	"net/http"

	"playit/consumer"
	"playit/messages"
	"playit/models"
	"playit/music"

	"github.com/labstack/echo/v4"
)

// RegisterAPIRoutes sets up the API routes
func RegisterAPIRoutes(e *echo.Echo, configs *models.Config) {
	e.GET("/api/healthz", healthCheck)
	e.GET("/api/debug/messages", getDebugMessages)
	e.GET("/api/debug/queue", getDebugQueue)

	e.GET("/auth", func(c echo.Context) error {
		authURL := fmt.Sprintf("https://id.twitch.tv/oauth2/authorize?response_type=code&client_id=%s&redirect_uri=%s&scope=chat:read+chat:edit",
			configs.TwitchClientID, configs.RedirectURI)
		return c.Redirect(http.StatusFound, authURL)
	})

	e.GET("/auth/callback", func(c echo.Context) error {
		return handleAuthCallback(c, configs)
	})
}

func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "UP"})
}

func getDebugMessages(c echo.Context) error {
	return c.JSON(http.StatusOK, messages.GetMessages())
}

func getDebugQueue(c echo.Context) error {
	return c.JSON(http.StatusOK, music.GetMusicQueue())
}

func handleAuthCallback(c echo.Context, configs *models.Config) error {
	code := c.QueryParam("code")
	if code == "" {
		return c.String(http.StatusBadRequest, "Code not provided")
	}

	// Exchange the authorization code for a user access token
	// token, err := consumer.ExchangeCodeForToken(code, configs.TwitchClientID, configs.TwitchClientSecret, configs.RedirectURI)
	// if err != nil {
	// 	log.Printf("Error exchanging code for token: %v\n", err)
	// 	return c.String(http.StatusInternalServerError, "Error retrieving access token")
	// }

	token := "weoali93356g04vm9534u61b05sxow"

	// Store the token globally
	configs.TwitchToken = token

	// Connect to Twitch chat using the updated token
	go consumer.ConnectAndConsumeTwitchChat(configs.TwitchChannelName, token)

	return c.JSON(http.StatusOK, map[string]string{"access_token": token})
}
