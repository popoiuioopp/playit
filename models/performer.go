package models

import "database/sql"

type Performer struct {
	Id             string         `json:"id"`
	Username       string         `json:"username"`
	YoutubeChannel sql.NullString `json:"youtube_channel"`
	TwitchChannel  sql.NullString `json:"twitch_channel"`
}
