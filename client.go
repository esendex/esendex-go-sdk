// Package xesende is a client for the Esendex REST API.
package xesende

import (
	"bytes"
	"encoding/xml"
	"net/http"
	"net/url"
)

const (
	defaultBaseURL   = "https://api.esendex.com/"
	defaultUserAgent = "xesende/golang"
)

// Client is the entry point for accessing the Esendex REST API.
type Client struct {
	client *http.Client
	user   string
	pass   string

	BaseURL   *url.URL
	UserAgent string
}

// New returns a new API client that authenticates with the credentials provided.
func New(user, pass string) *Client {
	baseURL, _ := url.Parse(defaultBaseURL)

	return &Client{
		client: http.DefaultClient,
		user:   user,
		pass:   pass,

		BaseURL:   baseURL,
		UserAgent: defaultUserAgent,
	}
}

func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	reqURL, err := c.BaseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if body != nil {
		if err := xml.NewEncoder(buf).Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, reqURL.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/xml")
	req.Header.Add("User-Agent", c.UserAgent)
	req.SetBasicAuth(c.user, c.pass)

	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
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

// AccountClient is a client scoped to a specific account reference.
type AccountClient struct {
	*Client
	reference string
}

// Account creates a client that can make requests scoped to a specific account
// reference.
func (c *Client) Account(reference string) *AccountClient {
	return &AccountClient{
		Client:    c,
		reference: reference,
	}
}
