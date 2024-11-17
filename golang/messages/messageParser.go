package messages

import (
	"playit/models"
	"playit/music"
	"strings"
)

var prefix = "!ขอเพลง"

// ParseSongRequest checks if the message is a song request based on the command prefix
func ParseSongRequest(msg Message) *models.SongRequest {
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

	return &models.SongRequest{
		Requester: msg.User,
		SongName:  songName,
		Artist:    artist,
	}
}

// HandleMessage processes incoming messages and checks for song requests
func HandleMessage(msg Message) {
	// AddMessage(msg)

	// Check if the message is a song request
	if songRequest := ParseSongRequest(msg); songRequest != nil {
		// Use the EnqueueSongRequest function instead of directly accessing the channel
		music.EnqueueSongRequest(*songRequest)
	}
}
