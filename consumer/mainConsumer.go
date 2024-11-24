package consumer

import (
	"playit/messages"
	"playit/models"
	"playit/repository"

	"github.com/gempir/go-twitch-irc/v4"
)

func StartUpConsumes(config models.Config) {
	performers, err := repository.GetPerformers()
	if err != nil {
		panic("Error getting performers")
	}
	client := twitch.NewAnonymousClient()

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		messages.HandleMessage("twitch", message.Channel, message)
	})

	for _, performer := range performers {

		if performer.TwitchChannel.Valid {
			client.Join(performer.TwitchChannel.String)
		}
	}
	err = client.Connect()
	if err != nil {
		panic(err)
	}
}
