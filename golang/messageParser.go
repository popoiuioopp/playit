package main

import (
	"strings"
)

var prefix = "!ขอเพลง"

// ParseSongRequest checks if the message is a song request based on the command prefix
func ParseSongRequest(msg Message) *SongRequest {
	if !strings.HasPrefix(msg.Content, prefix) {
		return nil
	}

	// Remove the command prefix and trim whitespace
	content := strings.TrimSpace(strings.TrimPrefix(msg.Content, prefix))

	// Split into song name and optional artist
	parts := strings.SplitN(content, " - ", 2)
	songName := strings.TrimSpace(parts[0])
	if songName == "" {
		return nil
	}

	artist := ""
	if len(parts) > 1 {
		artist = strings.TrimSpace(parts[1])
	}

	return &SongRequest{
		Requester: msg.User,
		SongName:  songName,
		Artist:    artist,
	}
}

// HandleMessage processes incoming messages and checks for song requests
func HandleMessage(msg Message) {
	AddMessage(msg)

	if songRequest := ParseSongRequest(msg); songRequest != nil {
		musicQueueChan <- *songRequest
	}
}
