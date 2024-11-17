package music

import (
	"log"
	"playit/realtime"
	"sync"
)

type SongRequest struct {
	Requester string `json:"requester"`
	SongName  string `json:"song_name"`
	Artist    string `json:"artist,omitempty"`
}

var musicQueue = struct {
	sync.RWMutex
	Queue []SongRequest
}{}

var musicQueueChan chan SongRequest

// Initialize the song request channel
func InitMusicQueue() {
	musicQueueChan = make(chan SongRequest, 100)
	go ProcessMusicQueue()
}

// AddSongRequest adds a song request directly to the queue
func AddSongRequest(song SongRequest) {
	musicQueue.Lock()
	defer musicQueue.Unlock()
	musicQueue.Queue = append(musicQueue.Queue, song)
	log.Printf("Added song request: %+v\n", song)

	realtime.BroadcastMessage(musicQueue.Queue)
}

// GetMusicQueue returns the current list of song requests
func GetMusicQueue() []SongRequest {
	musicQueue.RLock()
	defer musicQueue.RUnlock()
	return musicQueue.Queue
}

// EnqueueSongRequest sends a song request to the channel
func EnqueueSongRequest(song SongRequest) {
	musicQueueChan <- song
}

// ProcessMusicQueue listens for song requests from the channel and adds them to the queue
func ProcessMusicQueue() {
	for song := range musicQueueChan {
		AddSongRequest(song)
	}
}
