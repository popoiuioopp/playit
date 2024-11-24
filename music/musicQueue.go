package music

import (
	"playit/models"
)

func GetInitialMusicQueue(userId string) (*[]models.SongRequest, error) {
	var queue []models.SongRequest

	// queue, err := repository.GetSongRequestsByUserName(userId
	// if err != nil {
	// 	return &queue, err
	// }

	return &queue, nil
}
