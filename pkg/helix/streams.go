package helix

import (
	"strconv"
	"time"
)

const (
	STREAMURL = "https://api.twitch.tv/helix/streams"
)

type GetStreamsOpts struct {
	Before    string
	After     string
	First     int
	GameId    string
	Language  string
	UserId    string
	UserLogin string
}

type streamResp struct {
	Data []streamData
}

type streamData struct {
	GameId      string `json:"game_id"`
	GameName    string `json:"game_name"`
	Id          string
	Language    string
	StartedAt   time.Time `json:"started_at"`
	Title       string
	UserName    string `json:"user_name"`
	UserLogin   string `json:"user_login"`
	UserId      string `json:"user_id"`
	ViewerCount int    `json:"viewer_count"`
}

func (c *Client) GetStreams(opts *GetStreamsOpts) (*streamResp, error) {
	params := map[string]string{
		"before":     opts.Before,
		"after":      opts.After,
		"first":      strconv.Itoa(opts.First),
		"game_id":    opts.GameId,
		"language":   opts.Language,
		"user_id":    opts.UserId,
		"user_login": opts.UserLogin,
	}
	resp, err := c.getWithParams(params, STREAMURL, &streamResp{})
	if err != nil {
		return nil, err
	}

	r := resp.Result().(*streamResp)
	return r, nil
}
