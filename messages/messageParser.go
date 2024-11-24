package messages

import (
	"log"
	"playit/models"
	"playit/realtime"
	"playit/repository"
	"strings"

	"github.com/gempir/go-twitch-irc/v4"
)

var prefix = "ขอเพลง"

// ParseSongRequest checks if the message is a song request based on the command prefix
func ParseSongRequest(msg twitch.PrivateMessage) *models.SongRequest {
	if !strings.HasPrefix(msg.Message, prefix) {
		return nil
	}

	// Remove the command prefix and trim whitespace
	content := strings.TrimSpace(strings.TrimPrefix(msg.Message, prefix))

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
		Requester: msg.User.Name,
		SongName:  songName,
		Artist:    artist,
	}
}

// HandleMessage processes incoming messages and checks for song requests
func HandleMessage(platform, channelID string, msg twitch.PrivateMessage) {
	if songRequest := ParseSongRequest(msg); songRequest != nil {
		// Insert the song request into the database
		err := repository.InsertSongRequest(channelID, platform, songRequest)
		if err != nil {
			log.Printf("Error inserting song request: %v\n", err)
			return
		}

		// Fetch the updated queue for the channel
		queue, err := repository.GetSongRequestsByChannelID(channelID)
		if err != nil {
			log.Printf("Error fetching song requests for channel %s: %v\n", channelID, err)
			return
		}

		// Broadcast the updated queue to subscribed performers
		realtime.BroadcastMessage(channelID, platform, queue)
	}
}
