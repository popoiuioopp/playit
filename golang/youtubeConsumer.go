package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	YtChat "github.com/abhinavxd/youtube-live-chat-downloader/v2"
)

func consumeYouTubeChat(videoURL string) {
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

func StartYouTubeChatListener(channelID string) {
	for {
		liveURL, err := fetchLiveVideoURL(channelID)
		if err != nil {
			log.Printf("Error fetching live video URL for channel %s: %v\n", channelID, err)
			time.Sleep(30 * time.Second) // Retry after 30 seconds
			continue
		}

		if liveURL == "" {
			log.Printf("Channel %s is not live\n", channelID)
			time.Sleep(5 * time.Minute) // Check again after 5 minutes
			continue
		}

		// Start consuming YouTube chat
		go consumeYouTubeChat(liveURL)

		// Wait for the stream to end
		time.Sleep(10 * time.Minute) // Check every 10 minutes
	}
}

func fetchLiveVideoURL(channelID string) (string, error) {
	if ytAPIKey == "" {
		return "", fmt.Errorf("YouTube API key is missing")
	}

	apiURL := fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?part=snippet&channelId=%s&eventType=live&type=video&key=%s", channelID, ytAPIKey)
	resp, err := http.Get(apiURL)
	if err != nil {
		return "", fmt.Errorf("failed to make API request: %v", err)
	}
	defer resp.Body.Close()

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	// Parse the JSON response
	var result struct {
		Items []struct {
			ID struct {
				VideoID string `json:"videoId"`
			} `json:"id"`
		} `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	// Check if there is a live video
	if len(result.Items) > 0 && result.Items[0].ID.VideoID != "" {
		videoID := result.Items[0].ID.VideoID
		return "https://www.youtube.com/watch?v=" + videoID, nil
	}

	return "", nil // No live video found
}
