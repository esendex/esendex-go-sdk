package xesende_test

import (
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	xesende "."
	"github.com/stretchr/testify/assert"
)

func TestMessagesSent(t *testing.T) {
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
		username    = "user"
	)

	var (
		lastStatusAt    = time.Date(2012, 1, 1, 12, 0, 5, 0, time.UTC)
		lastStatusAtStr = "2012-01-01T12:00:05.000"
		submittedAt     = time.Date(2012, 1, 1, 12, 0, 2, 0, time.UTC)
		submittedAtStr  = "2012-01-01T12:00:02.000"
	)

	h := newRecordingHandler(`<?xml version="1.0" encoding="utf-8"?>
<messageheaders startindex="`+strconv.Itoa(startIndex)+`" count="`+strconv.Itoa(count)+`" totalcount="`+strconv.Itoa(totalCount)+`" xmlns="http://api.esendex.com/ns/">
 <messageheader id="`+id+`" uri="`+uri+`">
  <reference>`+reference+`</reference>
  <status>`+status+`</status>
  <laststatusat>`+lastStatusAtStr+`</laststatusat>
  <submittedat>`+submittedAtStr+`</submittedat>
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
  <username>`+username+`</username>
 </messageheader>
</messageheaders>`, 200, map[string]string{})
	s := httptest.NewServer(h)
	defer s.Close()

	client := xesende.New("user", "pass")
	client.BaseUrl, _ = url.Parse(s.URL)

	result, err := client.Messages.Sent()

	assert := assert.New(t)

	assert.Nil(err)

	assert.Equal("GET", h.Request.Method)
	assert.Equal("/v1.0/messageheaders", h.Request.URL.String())

	assert.Equal(startIndex, result.StartIndex)
	assert.Equal(count, result.Count)
	assert.Equal(totalCount, result.TotalCount)

	if assert.Equal(1, len(result.Messages)) {
		message := result.Messages[0]

		assert.Equal(id, message.Id)
		assert.Equal(uri, message.Uri)
		assert.Equal(reference, message.Reference)
		assert.Equal(status, message.Status)
		assert.Equal(lastStatusAt, message.LastStatusAt)
		assert.Equal(submittedAt, message.SubmittedAt)
		assert.Equal(messageType, message.Type)
		assert.Equal(to, message.To)
		assert.Equal(from, message.From)
		assert.Equal(summary, message.Summary)
		assert.Equal(bodyUri, message.BodyUri)
		assert.Equal(direction, message.Direction)
		assert.Equal(parts, message.Parts)
		assert.Equal(username, message.Username)
	}
}

func TestMessagesSentWithPaging(t *testing.T) {
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
		username    = "user"
	)

	var (
		lastStatusAt    = time.Date(2012, 1, 1, 12, 0, 5, 0, time.UTC)
		lastStatusAtStr = "2012-01-01T12:00:05.000"
		submittedAt     = time.Date(2012, 1, 1, 12, 0, 2, 0, time.UTC)
		submittedAtStr  = "2012-01-01T12:00:02.000"
	)

	h := newRecordingHandler(`<?xml version="1.0" encoding="utf-8"?>
<messageheaders startindex="`+strconv.Itoa(startIndex)+`" count="`+strconv.Itoa(count)+`" totalcount="`+strconv.Itoa(totalCount)+`" xmlns="http://api.esendex.com/ns/">
 <messageheader id="`+id+`" uri="`+uri+`">
  <reference>`+reference+`</reference>
  <status>`+status+`</status>
  <laststatusat>`+lastStatusAtStr+`</laststatusat>
  <submittedat>`+submittedAtStr+`</submittedat>
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
  <username>`+username+`</username>
 </messageheader>
</messageheaders>`, 200, map[string]string{})
	s := httptest.NewServer(h)
	defer s.Close()

	client := xesende.New("user", "pass")
	client.BaseUrl, _ = url.Parse(s.URL)

	result, err := client.Messages.Sent(xesende.Page(5, 10))

	assert := assert.New(t)

	assert.Nil(err)

	assert.Equal("GET", h.Request.Method)
	assert.Equal("/v1.0/messageheaders", h.Request.URL.Path)

	query := h.Request.URL.Query()
	assert.Equal("5", query.Get("startindex"))
	assert.Equal("10", query.Get("count"))

	assert.Equal(startIndex, result.StartIndex)
	assert.Equal(count, result.Count)
	assert.Equal(totalCount, result.TotalCount)

	if assert.Equal(1, len(result.Messages)) {
		message := result.Messages[0]

		assert.Equal(id, message.Id)
		assert.Equal(uri, message.Uri)
		assert.Equal(reference, message.Reference)
		assert.Equal(status, message.Status)
		assert.Equal(lastStatusAt, message.LastStatusAt)
		assert.Equal(submittedAt, message.SubmittedAt)
		assert.Equal(messageType, message.Type)
		assert.Equal(to, message.To)
		assert.Equal(from, message.From)
		assert.Equal(summary, message.Summary)
		assert.Equal(bodyUri, message.BodyUri)
		assert.Equal(direction, message.Direction)
		assert.Equal(parts, message.Parts)
		assert.Equal(username, message.Username)
	}
}
