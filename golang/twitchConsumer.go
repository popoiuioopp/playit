package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

func exchangeCodeForToken(code string) (string, error) {
	data := url.Values{}
	data.Set("client_id", twitchClientID)
	data.Set("client_secret", twitchClientSecret)
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

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}

	return tokenResp.AccessToken, nil
}

func connectAndConsumeTwitchChat(channelName string) {
	conn, _, err := websocket.DefaultDialer.Dial("wss://irc-ws.chat.twitch.tv:443", nil)
	if err != nil {
		log.Fatalf("Error connecting to Twitch IRC: %v\n", err)
	}
	defer conn.Close()

	// Authenticate with PASS and NICK commands
	if err := conn.WriteMessage(websocket.TextMessage, []byte("PASS oauth:"+config.Token)); err != nil {
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
			parsedMessage := parseMessage(string(message), channel)
			if parsedMessage != nil {
				HandleMessage(*parsedMessage)
			}
		}
	}
}
