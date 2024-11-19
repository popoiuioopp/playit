package models

import "time"

type Message struct {
	User      string    `json:"user"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
	Channel   string    `json:"channel"`
}
