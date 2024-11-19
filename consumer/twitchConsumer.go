package consumer

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"playit/messages"
	"playit/models"

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
	conn, _, err := websocket.DefaultDialer.Dial("wss://irc-ws.chat.twitch.tv:443", nil)
	if err != nil {
		log.Fatalf("Error connecting to Twitch IRC: %v\n", err)
	}
	defer conn.Close()

	// Authenticate using the token
	if err := conn.WriteMessage(websocket.TextMessage, []byte("PASS oauth:"+token)); err != nil {
		log.Println("Auth error:", err)
		return
	}
	if err := conn.WriteMessage(websocket.TextMessage, []byte("NICK kbj_bot")); err != nil {
		log.Println("Nick error:", err)
		return
	}

	// Request capabilities
	capReq := "CAP REQ :twitch.tv/membership twitch.tv/tags twitch.tv/commands"
	if err := conn.WriteMessage(websocket.TextMessage, []byte(capReq)); err != nil {
		log.Println("CAP REQ error:", err)
		return
	}

	// Join the specified channel
	if err := conn.WriteMessage(websocket.TextMessage, []byte("JOIN #"+channelName)); err != nil {
		log.Println("Join error:", err)
		return
	}

	log.Printf("Connected to %s chat\n", channelName)

	// Listen for messages from Twitch
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			return
		}

		// Respond to PING messages to keep the connection alive
		if string(message) == "PING :tmi.twitch.tv" {
			if err := conn.WriteMessage(websocket.TextMessage, []byte("PONG :tmi.twitch.tv")); err != nil {
				log.Println("PONG error:", err)
				return
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
