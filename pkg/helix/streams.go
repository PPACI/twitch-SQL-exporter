package helix

import "time"

type streamResp struct {
	Data []streamData
}

type streamData struct {
	GameId      string
	GameName    string
	Id          string
	Language    string
	StartedAt   time.Time
	Title       string
	UserName    string
	UserLogin   string
	UserId      string
	ViewerCount int
}
