package xesende

import "errors"

// AccountMessagesClient is a client scoped to a specific account for making
// requests for messages.
type AccountMessagesClient struct {
	*AccountClient
}

// Received returns the messages sent to the account.
func (c *AccountMessagesClient) Received(opts ...Option) (*MessagesReceivedResponse, error) {
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
			Id:         message.Id,
			Uri:        message.Uri,
			Reference:  message.Reference,
			Status:     message.Status,
			ReceivedAt: message.ReceivedAt.Time,
			Type:       message.Type,
			To:         message.To,
			From:       message.From,
			Summary:    message.Summary,
			BodyUri:    message.Body.Uri,
			Direction:  message.Direction,
			Parts:      message.Parts,
			ReadAt:     message.ReadAt.Time,
			ReadBy:     message.ReadBy,
		}
	}

	return response, nil
}
