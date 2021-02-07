package helix

import (
	"context"
	"errors"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/oauth2/twitch"
)

const (
	STREAMURL = "https://api.twitch.tv/helix/streams"
)

type ClientOpts struct {
	ClientID     string
	ClientSecret string
}

type Client struct {
	client *resty.Client
}

func NewClient(opts *ClientOpts, ctx context.Context) *Client {
	restyClient := newRestyClient(opts, ctx)
	return &Client{client: restyClient}
}

func newRestyClient(opts *ClientOpts, ctx context.Context) *resty.Client {
	oauth2Config := &clientcredentials.Config{
		ClientID:     opts.ClientID,
		ClientSecret: opts.ClientSecret,
		TokenURL:     twitch.Endpoint.TokenURL,
	}
	oauthClient := oauth2Config.Client(ctx)

	restyClient := resty.New()
	restyClient.SetHeader("Client-Id", opts.ClientID)
	restyClient.SetTransport(oauthClient.Transport)
	return restyClient
}

func (c *Client) GetStreams() ([]streamData, error) {
	resp, err := c.client.R().SetResult(&streamResp{}).Get(STREAMURL)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() > 399 {
		return nil, errors.New(string(resp.Body()))
	}
	data := resp.Result().(*streamResp)
	log.Debugln(resp)

	return data.Data, nil
}
