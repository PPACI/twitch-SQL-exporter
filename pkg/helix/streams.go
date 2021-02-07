package helix

import (
	"time"
)

const (
	STREAMURL = "https://api.twitch.tv/helix/streams"
)

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

func (c *Client) GetStreams() (*streamResp, error) {
	resp, err := c.client.R().
		SetResult(&streamResp{}).
		Get(STREAMURL)

	if err != nil {
		return nil, err
	}

	if e := resp.Error(); e != nil {
		return nil, e.(*helixError)
	}

	r := resp.Result().(*streamResp)

	return r, nil
}
