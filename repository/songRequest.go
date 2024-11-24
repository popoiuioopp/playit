package repository

import (
	"playit/db"
	"playit/models"
)

func InsertSongRequest(channelID, platform string, songRequest *models.SongRequest) error {
	query := `
		INSERT INTO request (
			channel_id,
			music,
			artist,
			status,
			requester,
			requester_platform
		) VALUES ($1, $2, $3, 'in_queue', $4, $5);`

	_, err := db.DB.Exec(query,
		channelID,
		songRequest.SongName,
		songRequest.Artist,
		songRequest.Requester,
		platform)
	if err != nil {
		return err
	}

	return nil
}

func GetSongRequestsByChannelID(channelID string) ([]models.SongRequest, error) {
	var songRequests []models.SongRequest
	query := `
		SELECT 
			music, 
			artist, 
			status, 
			requester, 
			requester_platform 
		FROM 
			request 
		WHERE 
			channel_id = $1;`

	rows, err := db.DB.Query(query, channelID)
	if err != nil {
		return songRequests, err
	}
	defer rows.Close()

	for rows.Next() {
		var songRequest models.SongRequest
		if err = rows.Scan(&songRequest.SongName, &songRequest.Artist, &songRequest.Status, &songRequest.Requester, &songRequest.RequesterPlatform); err != nil {
			return songRequests, err
		}
		songRequests = append(songRequests, songRequest)
	}

	return songRequests, nil
}

func GetPerformersByChannelAndPlatform(channelID, platform string) ([]models.Performer, error) {
	var performers []models.Performer

	query := `
		SELECT id, username, youtube_channel, twitch_channel 
		FROM performer 
		WHERE ($1 = 'twitch' AND twitch_channel = $2)
		   OR ($1 = 'youtube' AND youtube_channel = $2);`

	rows, err := db.DB.Query(query, platform, channelID)
	if err != nil {
		return performers, err
	}
	defer rows.Close()

	for rows.Next() {
		var performer models.Performer
		if err := rows.Scan(&performer.Id, &performer.Username, &performer.YoutubeChannel, &performer.TwitchChannel); err != nil {
			return performers, err
		}
		performers = append(performers, performer)
	}

	return performers, nil
}

func GetSongRequestsByUserName(userName string) ([]models.SongRequest, error) {
	var songRequests []models.SongRequest

	query := `
		SELECT 
			r.music, 
			r.artist, 
			r.status, 
			r.requester, 
			r.requester_platform 
		FROM 
			request r
		INNER JOIN 
			performer p
		ON 
			r.channel_id = p.twitch_channel OR r.channel_id = p.youtube_channel
		WHERE 
			p.username = $1;
		`

	rows, err := db.DB.Query(query, userName)
	if err != nil {
		return songRequests, err
	}
	defer rows.Close()

	for rows.Next() {
		var songRequest models.SongRequest
		if err = rows.Scan(&songRequest.SongName, &songRequest.Artist, &songRequest.Status, &songRequest.Requester, &songRequest.RequesterPlatform); err != nil {
			return songRequests, err
		}
		songRequests = append(songRequests, songRequest)
	}

	return songRequests, nil
}
