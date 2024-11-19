package models

type SongStatus string

const (
	StatusInQueue SongStatus = "in_queue"
	StatusPlaying SongStatus = "playing"
	StatusPlayed  SongStatus = "played"
)

type SongRequest struct {
	ID            int        `json:"id"`               // Unique ID for each song in the queue
	Requester     string     `json:"requester"`        // Who requested the song
	SongName      string     `json:"song_name"`        // Name of the song
	Artist        string     `json:"artist,omitempty"` // Artist of the song
	Status        SongStatus `json:"status"`           // Status of the song
	PerformanceID string     `json:"performance"`      // ID for the performance
}
