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

func (c *Client) requestWithParams(params map[string]string, outputType interface{}) *resty.Request {
	cleanedParams := getCleanMap(params)
	return c.client.R().SetQueryParams(cleanedParams).SetResult(outputType)
}

func (c *Client) getWithParams(params map[string]string, url string, outputType interface{}) (*resty.Response, error) {
	req := c.requestWithParams(params, outputType)
	res, err := req.Get(url)
	if err != nil {
		return &resty.Response{}, err
	}

	if e := res.Error(); e != nil {
		return &resty.Response{}, e.(*helixError)
	}
	return res, nil
}
