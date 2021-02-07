package helix

import (
	"context"
	"github.com/go-resty/resty/v2"
	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/oauth2/twitch"
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
	restyClient.SetError(&helixError{})
	return restyClient
}
