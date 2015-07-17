package xesende

import (
	"bytes"
	"encoding/xml"
	"net/http"
	"net/url"
)

const (
	defaultBaseUrl   = "https://api.esendex.com/"
	defaultUserAgent = ""
)

type Client struct {
	client    *http.Client
	BaseUrl   *url.URL
	UserAgent string

	Messages *MessagesClient
}

func New(user, pass string) *Client {
	baseUrl, _ := url.Parse(defaultBaseUrl)

	c := &Client{
		client:    http.DefaultClient,
		BaseUrl:   baseUrl,
		UserAgent: defaultUserAgent,
	}

	c.Messages = &MessagesClient{c}

	return c
}

func (c *Client) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	reqUrl, err := c.BaseUrl.Parse(path)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if body != nil {
		if err := xml.NewEncoder(buf).Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, reqUrl.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/xml")
	req.Header.Add("User-Agent", c.UserAgent)

	return req, nil
}

func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		if rerr := resp.Body.Close(); err == nil {
			err = rerr
		}
	}()

	if v != nil {
		if err := xml.NewDecoder(resp.Body).Decode(v); err != nil {
			return nil, err
		}
	}

	return resp, err
}

type AccountClient struct {
	*Client
	reference string
}

func (c *Client) Account(reference string) *AccountClient {
	return &AccountClient{c, reference}
}
