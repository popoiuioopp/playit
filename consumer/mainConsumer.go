package consumer

import (
	"playit/messages"
	"playit/models"
	"playit/realtime"
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
			realtime.SubscribePerformerToChannel(performer.Username, "twitch", performer.TwitchChannel.String)
			client.Join(performer.TwitchChannel.String)
		}

		if performer.YoutubeChannel.Valid {
			realtime.SubscribePerformerToChannel(performer.Username, "youtube", performer.YoutubeChannel.String)
			// TODO: add yt consumer
		}
	}
	err = client.Connect()
	if err != nil {
		panic(err)
	}
}
