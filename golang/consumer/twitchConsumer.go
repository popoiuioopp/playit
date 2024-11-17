package consumer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"playit/messages"

	"github.com/gorilla/websocket"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// GetTwitchCredential retrieves an app access token using the client credentials grant flow.
func GetTwitchCredential(clientId, clientSecret string) (string, error) {
	apiURL := "https://id.twitch.tv/oauth2/token"

	// Prepare the form data
	data := url.Values{}
	data.Set("client_id", clientId)
	data.Set("client_secret", clientSecret)
	data.Set("grant_type", "client_credentials")

	// Create a POST request
	req, err := http.NewRequest("POST", apiURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make API request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	// Parse the response
	var result TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	return result.AccessToken, nil
}

func ConnectAndConsumeTwitchChat(channelName, twitchToken string) {
	conn, _, err := websocket.DefaultDialer.Dial("wss://irc-ws.chat.twitch.tv:443", nil)
	if err != nil {
		log.Fatalf("Error connecting to Twitch IRC: %v\n", err)
	}
	defer conn.Close()

	// Authenticate with PASS and NICK commands
	if err := conn.WriteMessage(websocket.TextMessage, []byte("PASS oauth:"+twitchToken)); err != nil {
		log.Println("Auth error:", err)
		return
	}
	if err := conn.WriteMessage(websocket.TextMessage, []byte("NICK your_twitch_username")); err != nil {
		log.Println("Nick error:", err)
		return
	}

	// Request capabilities
	capReq := "CAP REQ :twitch.tv/membership twitch.tv/tags twitch.tv/commands"
	if err := conn.WriteMessage(websocket.TextMessage, []byte(capReq)); err != nil {
		log.Println("CAP REQ error:", err)
		return
	}

	// Join a chat room
	channel := channelName
	if err := conn.WriteMessage(websocket.TextMessage, []byte("JOIN "+channel)); err != nil {
		log.Println("Join error:", err)
		return
	}

	log.Printf("Connected to %s chat\n", channel)

	// Listen for messages from Twitch
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			return
		}

		if string(message) == "PING :tmi.twitch.tv" {
			if err := conn.WriteMessage(websocket.TextMessage, []byte("PONG :tmi.twitch.tv")); err != nil {
				log.Println("PONG error:", err)
				return
			}
		} else {
			parsedMessage := messages.ParseMessage(string(message), channel)
			if parsedMessage != nil {
				messages.HandleMessage(*parsedMessage)
			}
		}
	}
}
