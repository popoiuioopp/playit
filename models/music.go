package models

type SongRequest struct {
	Requester string `json:"requester"`
	SongName  string `json:"song_name"`
	Artist    string `json:"artist,omitempty"`
}
