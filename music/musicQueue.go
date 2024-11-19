package music

import (
	"log"
	"playit/models"
	"playit/realtime"
	"sync"
	"time"
)

var musicQueue = struct {
	sync.RWMutex
	Queues map[string][]models.SongRequest
}{
	Queues: make(map[string][]models.SongRequest),
}

var musicQueueChan chan models.SongRequest

// Initialize the song request channel
func InitMusicQueue() {
	musicQueueChan = make(chan models.SongRequest, 100)
	go ProcessMusicQueue()
}

// AddSongRequest adds a song request directly to the queue
func AddSongRequest(song models.SongRequest) {
	musicQueue.Lock()
	defer musicQueue.Unlock()

	// Generate a unique ID for the song
	song.ID = len(musicQueue.Queues[song.PerformanceID]) + 1
	song.Status = models.StatusInQueue

	musicQueue.Queues[song.PerformanceID] = append(musicQueue.Queues[song.PerformanceID], song)
	log.Printf("Added song request: %+v\n", song)
	realtime.BroadcastMessage(musicQueue.Queues[song.PerformanceID])
}

func SkipSong(performanceID string) {
	musicQueue.Lock()
	defer musicQueue.Unlock()

	if len(musicQueue.Queues[performanceID]) == 0 {
		log.Println("No songs to skip")
		return
	}

	// Remove the first song from the queue
	musicQueue.Queues[performanceID] = musicQueue.Queues[performanceID][1:]
	realtime.BroadcastMessage(musicQueue.Queues[performanceID])
}

func MarkSongAsPlayed(performanceID string) {
	musicQueue.Lock()
	defer musicQueue.Unlock()

	if len(musicQueue.Queues[performanceID]) == 0 {
		log.Println("No songs to mark as played")
		return
	}

	// Mark the first song as played
	musicQueue.Queues[performanceID][0].Status = models.StatusPlayed

	// Move to the next song and mark it as playing, if available
	if len(musicQueue.Queues[performanceID]) > 1 {
		musicQueue.Queues[performanceID][1].Status = models.StatusPlaying
	}
	realtime.BroadcastMessage(musicQueue.Queues[performanceID])
}

func DeleteSong(performanceID string, songID int) {
	musicQueue.Lock()
	defer musicQueue.Unlock()

	queue := musicQueue.Queues[performanceID]
	for i, song := range queue {
		if song.ID == songID {
			// Remove the song from the queue
			musicQueue.Queues[performanceID] = append(queue[:i], queue[i+1:]...)
			log.Printf("Deleted song with ID: %d\n", songID)
			break
		}
	}
	realtime.BroadcastMessage(musicQueue.Queues[performanceID])
}

// GetMusicQueue returns the current list of song requests
func GetMusicQueue(performanceID string) []models.SongRequest {
	musicQueue.RLock()
	defer musicQueue.RUnlock()
	return musicQueue.Queues[performanceID]
}

// EnqueueSongRequest sends a song request to the channel
func EnqueueSongRequest(song models.SongRequest) {
	musicQueueChan <- song
}

// ProcessMusicQueue listens for song requests from the channel and adds them to the queue
func ProcessMusicQueue() {
	for song := range musicQueueChan {
		AddSongRequest(song)
	}
}

func GetTodayPerformanceID() string {
	return time.Now().Format("2024-11-17") // e.g., "2024-11-17"
}
