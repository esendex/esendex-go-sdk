package xesende_test

import (
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"hawx.me/code/xesende"
)

func TestAccountMessagesReceived(t *testing.T) {
	const (
		startIndex  = 0
		count       = 15
		totalCount  = 200
		id          = "messageheaderid"
		uri         = "http://somemessageheader"
		reference   = "EXETRTRE"
		status      = "STATUS"
		messageType = "TYPE"
		to          = "4538224364236"
		from        = "428377843"
		summary     = "SUM"
		bodyUri     = "http://rrehekr"
		direction   = "OUT"
		parts       = 1
		readBy      = "someone"
	)

	var (
		receivedAt    = time.Date(2012, 1, 1, 12, 0, 5, 0, time.UTC)
		receivedAtStr = "2012-01-01T12:00:05"
		readAt        = time.Date(2012, 1, 1, 12, 0, 2, 0, time.UTC)
		readAtStr     = "2012-01-01T12:00:02"
	)

	h := newRecordingHandler(`<?xml version="1.0" encoding="utf-8"?>
<messageheaders startindex="`+strconv.Itoa(startIndex)+`" count="`+strconv.Itoa(count)+`" totalcount="`+strconv.Itoa(totalCount)+`" xmlns="http://api.esendex.com/ns/">
 <messageheader id="`+id+`" uri="`+uri+`">
  <reference>`+reference+`</reference>
  <status>`+status+`</status>
  <receivedat>`+receivedAtStr+`</receivedat>
  <type>`+messageType+`</type>
  <to>
    <phonenumber>`+to+`</phonenumber>
  </to>
  <from>
   <phonenumber>`+from+`</phonenumber>
  </from>
  <summary>`+summary+`</summary>
  <body uri="`+bodyUri+`"/>
  <direction>`+direction+`</direction>
  <parts>`+strconv.Itoa(parts)+`</parts>
  <readat>`+readAtStr+`</readat>
  <readby>`+readBy+`</readby>
 </messageheader>
</messageheaders>`, 200, map[string]string{})
	s := httptest.NewServer(h)
	defer s.Close()

	client := xesende.New("user", "pass")
	client.BaseUrl, _ = url.Parse(s.URL)

	account := client.Account("EXHEY")
	result, err := account.Messages.Received()

	assert := assert.New(t)

	assert.Nil(err)

	assert.Equal("GET", h.Request.Method)
	assert.Equal("/v1.0/inbox/EXHEY/messages", h.Request.URL.String())

	if user, pass, ok := h.Request.BasicAuth(); assert.True(ok) {
		assert.Equal("user", user)
		assert.Equal("pass", pass)
	}

	assert.Equal(startIndex, result.StartIndex)
	assert.Equal(count, result.Count)
	assert.Equal(totalCount, result.TotalCount)

	if assert.Equal(1, len(result.Messages)) {
		message := result.Messages[0]

		assert.Equal(id, message.Id)
		assert.Equal(uri, message.Uri)
		assert.Equal(reference, message.Reference)
		assert.Equal(status, message.Status)
		assert.Equal(receivedAt, message.ReceivedAt)
		assert.Equal(messageType, message.Type)
		assert.Equal(to, message.To)
		assert.Equal(from, message.From)
		assert.Equal(summary, message.Summary)
		assert.Equal(bodyUri, message.BodyUri)
		assert.Equal(direction, message.Direction)
		assert.Equal(parts, message.Parts)
		assert.Equal(readAt, message.ReadAt)
		assert.Equal(readBy, message.ReadBy)
	}
}
