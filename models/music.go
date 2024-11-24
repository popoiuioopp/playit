package models

type SongRequest struct {
	Requester         string `json:"requester"`          // Who requested the song
	SongName          string `json:"song_name"`          // Name of the song
	Artist            string `json:"artist,omitempty"`   // Artist of the song
	Status            string `json:"status"`             // Status of the song
	RequesterPlatform string `json:"requester_platform"` // Platform of the requester (eg. youtube, twitch)
}

type Request struct {
	ID                string `json:"id"`
	PerformanceId     string `json:"performance_id"`
	Music             string `json:"music"`
	Artist            string `json:"artist"`
	Status            string `json:"status"`
	Requester         string `json:"requester"`
	RequesterPlatform string `json:"requester_platform"`
}
