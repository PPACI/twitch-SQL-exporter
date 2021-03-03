package sql

import (
	"gorm.io/gorm"
	"time"
)

type StreamRecord struct {
	gorm.Model

	StreamId      string `gorm:"index"`
	GameId        string
	GameName      string `gorm:"index"`
	Language      string
	StartedAt     time.Time
	Title         string
	UserName      string `gorm:"index"`
	UserLogin     string
	UserId        string `gorm:"index"`
	ViewerCount   int
	FollowerCount int
}
