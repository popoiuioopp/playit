package consumer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"playit/messages"
	"playit/models"
	"time"

	"github.com/gorilla/websocket"
)

func ExchangeCodeForToken(code, clientID, clientSecret, redirectURI string) (string, error) {
	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", redirectURI)
	req, err := http.NewRequest("POST", "https://id.twitch.tv/oauth2/token", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var tokenResp models.TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}
	return tokenResp.AccessToken, nil
}

func ConnectAndConsumeTwitchChat(channelName, token string) {
	for {
		err := connectToTwitchChat(channelName, token)
		if err != nil {
			log.Printf("Error in Twitch chat connection: %v. Reconnecting...\n", err)
		}

		// Wait for a while before attempting to reconnect
		time.Sleep(5 * time.Second)
	}
}

func connectToTwitchChat(channelName, token string) error {
	conn, _, err := websocket.DefaultDialer.Dial("wss://irc-ws.chat.twitch.tv:443", nil)
	if err != nil {
		return fmt.Errorf("error connecting to Twitch IRC: %v", err)
	}
	defer conn.Close()

	// Authenticate using the token
	if err := conn.WriteMessage(websocket.TextMessage, []byte("PASS oauth:"+token)); err != nil {
		return fmt.Errorf("auth error: %v", err)
	}
	if err := conn.WriteMessage(websocket.TextMessage, []byte("NICK kbj_bot")); err != nil {
		return fmt.Errorf("nick error: %v", err)
	}

	// Request capabilities
	capReq := "CAP REQ :twitch.tv/membership twitch.tv/tags twitch.tv/commands"
	if err := conn.WriteMessage(websocket.TextMessage, []byte(capReq)); err != nil {
		return fmt.Errorf("CAP REQ error: %v", err)
	}

	// Join the specified channel
	if err := conn.WriteMessage(websocket.TextMessage, []byte("JOIN #"+channelName)); err != nil {
		return fmt.Errorf("join error: %v", err)
	}

	log.Printf("Connected to %s chat\n", channelName)

	// Listen for messages from Twitch
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("read error: %v", err)
			continue
		}

		// Respond to PING messages to keep the connection alive
		if string(message) == "PING :tmi.twitch.tv" {
			if err := conn.WriteMessage(websocket.TextMessage, []byte("PONG :tmi.twitch.tv")); err != nil {
				return fmt.Errorf("PONG error: %v", err)
			}
		} else {
			parsedMessage := messages.ParseMessage(string(message), channelName)
			if parsedMessage != nil {
				messages.HandleMessage(*parsedMessage)
				log.Printf("Message received: %v\n", parsedMessage)
			}
		}
	}
}
