package xesende

import (
	"encoding/xml"
	"errors"
	"time"
)

type Paging struct {
	StartIndex int
	Count      int
	TotalCount int
}

type MessagesResponse struct {
	Paging
	Messages []MessageResponse
}

type MessageResponse struct {
	Id           string
	Uri          string
	Reference    string
	Status       string
	LastStatusAt time.Time
	SubmittedAt  time.Time
	Type         string
	To           string
	From         string
	Summary      string
	BodyUri      string
	Direction    string
	Parts        int
	Username     string
}
type MessagesReceivedResponse struct {
	Paging
	Messages []MessageReceivedResponse
}

type MessageReceivedResponse struct {
	Id         string
	Uri        string
	Reference  string
	Status     string
	ReceivedAt time.Time
	Type       string
	To         string
	From       string
	Summary    string
	BodyUri    string
	Direction  string
	Parts      int
	ReadAt     time.Time
	ReadBy     string
}

type MessagesClient struct {
	*Client
}

func (c *MessagesClient) Sent(opts ...Option) (*MessagesResponse, error) {
	req, err := c.NewRequest("GET", "/v1.0/messageheaders", nil)
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		opt(req)
	}

	var v messageHeadersResponse
	resp, err := c.Do(req, &v)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("Expected 200")
	}

	response := &MessagesResponse{
		Paging: Paging{
			StartIndex: v.StartIndex,
			Count:      v.Count,
			TotalCount: v.TotalCount,
		},
		Messages: make([]MessageResponse, len(v.Messages)),
	}

	for i, message := range v.Messages {
		response.Messages[i] = MessageResponse{
			Id:           message.Id,
			Uri:          message.Uri,
			Reference:    message.Reference,
			Status:       message.Status,
			LastStatusAt: message.LastStatusAt.Time,
			SubmittedAt:  message.SubmittedAt.Time,
			Type:         message.Type,
			To:           message.To,
			From:         message.From,
			Summary:      message.Summary,
			BodyUri:      message.Body.Uri,
			Direction:    message.Direction,
			Parts:        message.Parts,
			Username:     message.Username,
		}
	}

	return response, nil
}

func (c *MessagesClient) Received(opts ...Option) (*MessagesReceivedResponse, error) {
	req, err := c.NewRequest("GET", "/v1.0/inbox/messages", nil)
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		opt(req)
	}

	var v inboxResponse
	resp, err := c.Do(req, &v)
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

func (c *MessagesClient) ById(id string) (*MessageResponse, error) {
	req, err := c.NewRequest("GET", "/v1.0/messageheaders/"+id, nil)
	if err != nil {
		return nil, err
	}

	var v messageHeadersResponseMessageHeader
	resp, err := c.Do(req, &v)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("Expected 200")
	}

	response := &MessageResponse{
		Id:           v.Id,
		Uri:          v.Uri,
		Reference:    v.Reference,
		Status:       v.Status,
		LastStatusAt: v.LastStatusAt.Time,
		SubmittedAt:  v.SubmittedAt.Time,
		Type:         v.Type,
		To:           v.To,
		From:         v.From,
		Summary:      v.Summary,
		BodyUri:      v.Body.Uri,
		Direction:    v.Direction,
		Parts:        v.Parts,
		Username:     v.Username,
	}

	return response, nil
}

type messageHeadersResponse struct {
	XMLName    xml.Name                              `xml:"http://api.esendex.com/ns/ messageheaders"`
	StartIndex int                                   `xml:"startindex,attr"`
	Count      int                                   `xml:"count,attr"`
	TotalCount int                                   `xml:"totalcount,attr"`
	Messages   []messageHeadersResponseMessageHeader `xml:"messageheader"`
}

type messageHeadersResponseMessageHeader struct {
	Id           string            `xml:"id,attr"`
	Uri          string            `xml:"uri,attr"`
	Reference    string            `xml:"reference"`
	Status       string            `xml:"status"`
	LastStatusAt messageHeaderTime `xml:"laststatusat"`
	SubmittedAt  messageHeaderTime `xml:"submittedat"`
	Type         string            `xml:"type"`
	To           string            `xml:"to>phonenumber"`
	From         string            `xml:"from>phonenumber"`
	Summary      string            `xml:"summary"`
	Body         struct {
		Uri string `xml:"uri,attr"`
	} `xml:"body"`
	Direction string `xml:"direction"`
	Parts     int    `xml:"parts"`
	Username  string `xml:"username"`
}

type inboxResponse struct {
	XMLName    xml.Name                     `xml:"http://api.esendex.com/ns/ messageheaders"`
	StartIndex int                          `xml:"startindex,attr"`
	Count      int                          `xml:"count,attr"`
	TotalCount int                          `xml:"totalcount,attr"`
	Messages   []inboxResponseMessageHeader `xml:"messageheader"`
}

type inboxResponseMessageHeader struct {
	Id         string            `xml:"id,attr"`
	Uri        string            `xml:"uri,attr"`
	Reference  string            `xml:"reference"`
	Status     string            `xml:"status"`
	ReceivedAt messageHeaderTime `xml:"receivedat"`
	Type       string            `xml:"type"`
	To         string            `xml:"to>phonenumber"`
	From       string            `xml:"from>phonenumber"`
	Summary    string            `xml:"summary"`
	Body       struct {
		Uri string `xml:"uri,attr"`
	} `xml:"body"`
	Direction string            `xml:"direction"`
	Parts     int               `xml:"parts"`
	ReadAt    messageHeaderTime `xml:"readat"`
	ReadBy    string            `xml:"readby"`
}

const messageHeaderTimeFormat = "2006-01-02T15:04:05.999999999"

type messageHeaderTime struct {
	time.Time
}

func (t messageHeaderTime) MarshalText() ([]byte, error) {
	return []byte(t.Format(messageHeaderTimeFormat)), nil
}

func (t *messageHeaderTime) UnmarshalText(data []byte) error {
	g, err := time.ParseInLocation(messageHeaderTimeFormat, string(data), time.UTC)
	if err != nil {
		return err
	}
	*t = messageHeaderTime{g}
	return nil
}
