package xesende_test

import (
	"net/http/httptest"
	"net/url"
	"testing"

	xesende "."
	"github.com/stretchr/testify/assert"
)

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

	assert.Equal(batchId, result.BatchId)

	assert.Equal(1, len(result.Messages))
	assert.Equal(messageId, result.Messages[0].Id)
	assert.Equal(messageUri, result.Messages[0].Uri)
}
