package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	YtChat "github.com/abhinavxd/youtube-live-chat-downloader/v2"
)

func StartYouTubeChatListener(videoURL string) {
	// Optional: Add custom cookies to bypass region-based consent pop-ups or CAPTCHA
	customCookies := []*http.Cookie{
		{Name: "PREF", Value: "tz=Europe.Rome", MaxAge: 300},
		{Name: "CONSENT", Value: fmt.Sprintf("YES+yt.432048971.it+FX+%d", 100+rand.Intn(999-100+1)), MaxAge: 300},
	}
	YtChat.AddCookies(customCookies)

	// Parse the initial data to get the continuation token
	continuation, cfg, err := YtChat.ParseInitialData(videoURL)
	if err != nil {
		log.Fatalf("Error parsing initial data: %v\n", err)
		return
	}

	fmt.Printf("Started listening to YouTube live chat for video: %s\n", videoURL)

	for {
		// Fetch chat messages using the continuation token
		chatMessages, newContinuation, err := YtChat.FetchContinuationChat(continuation, cfg)
		if err == YtChat.ErrLiveStreamOver {
			log.Println("Live stream is over")
			return
		}
		if err != nil {
			log.Printf("Error fetching chat messages: %v\n", err)
			continue
		}

		// Update the continuation token for the next request
		continuation = newContinuation

		// Process each chat message
		for _, msg := range chatMessages {
			// Create a new message struct and add it to the store
			parsedMessage := Message{
				User:      msg.AuthorName,
				Content:   msg.Message,
				Timestamp: time.Now(),
				Channel:   videoURL,
			}
			HandleMessage(parsedMessage)
		}

		time.Sleep(2 * time.Second)
	}
}
