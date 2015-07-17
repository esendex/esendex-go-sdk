package xesende_test

import (
	"net/http/httptest"
	"net/url"
	"testing"

	xesende "."
	"github.com/stretchr/testify/assert"
)

func ExampleAccountClient_Send() {
	accountClient := xesende.New("user@example.com", "pass").Account("EX00000")

	accountClient.Send(xesende.Messages{
		{To: "00000000", Body: "Hello"},
	})
}

func TestSendSingleMessage(t *testing.T) {
	const (
		batchId    = "batchId"
		messageId  = "messageId"
		messageUri = "messageUri"
	)

	h := newRecordingHandler(`<?xml version="1.0" encoding="utf-8"?>
<messageheaders batchid="`+batchId+`" xmlns="http://api.esendex.com/ns/">
  <messageheader uri="`+messageUri+`" id="`+messageId+`" />
</messageheaders>`, 200, map[string]string{})
	s := httptest.NewServer(h)
	defer s.Close()

	client := xesende.New("user", "pass")
	client.BaseUrl, _ = url.Parse(s.URL)

	account := client.Account("EXWHATEVS")

	result, err := account.Send(xesende.Messages{
		{To: "358973", Body: "HWEYERW"},
	})

	assert := assert.New(t)

	assert.Nil(err)

	assert.Equal("POST", h.Request.Method)
	assert.Equal("/v1.0/messagedispatcher", h.Request.URL.String())

	if user, pass, ok := h.Request.BasicAuth(); assert.True(ok) {
		assert.Equal("user", user)
		assert.Equal("pass", pass)
	}

	assert.Equal(batchId, result.BatchId)

	assert.Equal(1, len(result.Messages))
	assert.Equal(messageId, result.Messages[0].Id)
	assert.Equal(messageUri, result.Messages[0].Uri)
}
