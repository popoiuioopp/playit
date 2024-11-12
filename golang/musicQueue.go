package main

import (
	"log"
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

// Go channel for processing song requests
var musicQueueChan = make(chan SongRequest, 100)

// AddSongRequest adds a song request to the music queue
func AddSongRequest(song SongRequest) {
	musicQueue.Lock()
	defer musicQueue.Unlock()
	musicQueue.Queue = append(musicQueue.Queue, song)
	log.Printf("Added song request: %+v\n", song)
}

// GetMusicQueue returns the current list of song requests
func GetMusicQueue() []SongRequest {
	musicQueue.RLock()
	defer musicQueue.RUnlock()
	return musicQueue.Queue
}

// ProcessMusicQueue listens for song requests from the channel and adds them to the queue
func ProcessMusicQueue() {
	for song := range musicQueueChan {
		AddSongRequest(song)
	}
}
