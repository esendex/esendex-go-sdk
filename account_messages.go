package xesende

import (
	"errors"
	"net/http"
)

// Sent returns a list of messages sent by the account.
func (c *AccountClient) Sent(opts ...Option) (*MessagesResponse, error) {
	accountOption := func(r *http.Request) {
		q := r.URL.Query()

		q.Add("accountReference", c.reference)

		r.URL.RawQuery = q.Encode()
	}

	return c.Client.Sent(append(opts, accountOption)...)
}

// Received returns the messages sent to the account.
func (c *AccountClient) Received(opts ...Option) (*MessagesReceivedResponse, error) {
	req, err := c.newRequest("GET", "/v1.0/inbox/"+c.reference+"/messages", nil)
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		opt(req)
	}

	var v inboxResponse
	resp, err := c.do(req, &v)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("Expected 200")
	}

	response := &MessagesReceivedResponse{
		Paging: Paging{
			StartIndex: v.StartIndex,
			Count:      v.Count,
			TotalCount: v.TotalCount,
		},
		Messages: make([]MessageReceivedResponse, len(v.Messages)),
	}

	for i, message := range v.Messages {
		response.Messages[i] = MessageReceivedResponse{
			ID:         message.ID,
			URI:        message.URI,
			Reference:  message.Reference,
			Status:     message.Status,
			ReceivedAt: message.ReceivedAt.Time,
			Type:       message.Type,
			To:         message.To,
			From:       message.From,
			Summary:    message.Summary,
			BodyURI:    message.Body.URI,
			Direction:  message.Direction,
			Parts:      message.Parts,
			ReadAt:     message.ReadAt.Time,
			ReadBy:     message.ReadBy,
		}
	}

	return response, nil
}
