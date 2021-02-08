package sql

import (
	"gorm.io/gorm"
	"time"
)

type StreamRecord struct {
	gorm.Model

	StreamId    string
	GameId      string
	GameName    string
	Language    string
	StartedAt   time.Time
	Title       string
	UserName    string
	UserLogin   string
	UserId      string
	ViewerCount int
}
