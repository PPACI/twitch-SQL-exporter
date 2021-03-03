package helix

const (
	FOLLOWURL = "https://api.twitch.tv/helix/users/follows"
)

type GetFollowsOpts struct {
	UserId string
}

type followResp struct {
	Total int
}

func (c *Client) GetFollows(opts *GetFollowsOpts) (*followResp, error) {
	params := map[string]string{
		"to_id": opts.UserId,
	}
	resp, err := c.getWithParams(params, FOLLOWURL, &followResp{})
	if err != nil {
		return nil, err
	}

	r := resp.Result().(*followResp)
	return r, nil
}
