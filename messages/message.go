package messages

import (
	"playit/models"
	"strings"
	"sync"
	"time"
)

var messageStore = struct {
	sync.RWMutex
	Messages []models.Message
}{}

func AddMessage(msg models.Message) {
	messageStore.Lock()
	defer messageStore.Unlock()
	messageStore.Messages = append(messageStore.Messages, msg)
}

func GetMessages() []models.Message {
	messageStore.RLock()
	defer messageStore.RUnlock()
	return messageStore.Messages
}

// parseMessage extracts the user and content from the raw IRC message received from Twitch chat.
//
// Example raw message:
// "@badge-info=subscriber/1;badges=subscriber/3003;display-name=Ka_BeeJa;mod=0;user-id=12345; :ka_beeja!ka_beeja@ka_beeja.tmi.twitch.tv PRIVMSG #ka_beeja :Hello World!"
//
// Extracted fields:
// - User: "Ka_BeeJa"
// - Content: "Hello World!"
// - Channel: "#ka_beeja"
// - Timestamp: The current time when the message was parsed.
func ParseMessage(rawMessage, channelName string) *models.Message {
	if strings.Contains(rawMessage, "PRIVMSG") {
		// Split the tags and the actual message
		parts := strings.SplitN(rawMessage, " PRIVMSG ", 2)
		if len(parts) < 2 {
			return nil
		}

		// Extract the content part after the colon (":")
		contentParts := strings.SplitN(parts[1], " :", 2)
		if len(contentParts) < 2 {
			return nil
		}
		content := strings.TrimSpace(contentParts[1])

		// Extract the username from tags or prefix
		username := extractDisplayName(parts[0])
		if username == "" {
			username = extractUsername(parts[0])
		}

		return &models.Message{
			User:      username,
			Content:   content,
			Timestamp: time.Now(),
			Channel:   channelName,
		}
	}
	return nil
}

// extractDisplayName retrieves the "display-name" from the tags if available
func extractDisplayName(tags string) string {
	for _, tag := range strings.Split(tags, ";") {
		if strings.HasPrefix(tag, "display-name=") {
			return strings.SplitN(tag, "=", 2)[1]
		}
	}
	return ""
}

// extractUsername retrieves the username from the prefix if "display-name" is not available
func extractUsername(tags string) string {
	prefixParts := strings.Split(tags, " ")
	if len(prefixParts) > 1 {
		userPart := prefixParts[0]
		userParts := strings.Split(userPart, "!")
		if len(userParts) > 0 {
			return strings.TrimPrefix(userParts[0], ":")
		}
	}
	return ""
}
