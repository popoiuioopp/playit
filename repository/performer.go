package repository

import (
	"playit/db"
	"playit/models"
)

func GetPerformers() ([]models.Performer, error) {
	var performers []models.Performer

	rows, err := db.DB.Query("select id, username, twitch_channel, youtube_channel from performer;")
	if err != nil {
		return performers, err
	}
	defer rows.Close()

	for rows.Next() {
		var performer models.Performer
		if err = rows.Scan(&performer.Id, &performer.Username, &performer.TwitchChannel, &performer.YoutubeChannel); err != nil {
			return performers, err
		}
		performers = append(performers, performer)
	}

	return performers, nil
}

func GetPerformersByChannelIdAndPlatform(channel, platform string) ([]string, error) {
	var performers []string

	rows, err := db.DB.Query("select username from performer where (twitch_channel = $1 and $2 = 'twitch')")
	if err != nil {
		return performers, err
	}
	defer rows.Close()

	for rows.Next() {
		var performerUserName string
		if err = rows.Scan(&performerUserName); err != nil {
			return performers, err
		}
		performers = append(performers, performerUserName)
	}

	return performers, nil
}
