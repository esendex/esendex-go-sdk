package xesende_test

import (
	"fmt"
	"log"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	xesende "."
	"github.com/stretchr/testify/assert"
)

func ExampleMessagesClient_Sent() {
	client := xesende.New("user@example.com", "pass")

	response, err := client.Messages.Sent()
	if err != nil {
		log.Fatal(err)
	}

	for _, message := range response.Messages {
		fmt.Printf("%v %s: %s\n", message.SubmittedAt, message.To, message.Summary)
	}
}

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

	if user, pass, ok := h.Request.BasicAuth(); assert.True(ok) {
		assert.Equal("user", user)
		assert.Equal("pass", pass)
	}

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

func TestMessagesById(t *testing.T) {
	const (
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
<messageheader id="`+id+`" uri="`+uri+`" xmlns="http://api.esendex.com/ns/">
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
</messageheader>`, 200, map[string]string{})
	s := httptest.NewServer(h)
	defer s.Close()

	client := xesende.New("user", "pass")
	client.BaseUrl, _ = url.Parse(s.URL)

	result, err := client.Messages.ById(id)

	assert := assert.New(t)

	assert.Nil(err)

	assert.Equal("GET", h.Request.Method)
	assert.Equal("/v1.0/messageheaders/"+id, h.Request.URL.String())

	assert.Equal(id, result.Id)
	assert.Equal(uri, result.Uri)
	assert.Equal(reference, result.Reference)
	assert.Equal(status, result.Status)
	assert.Equal(lastStatusAt, result.LastStatusAt)
	assert.Equal(submittedAt, result.SubmittedAt)
	assert.Equal(messageType, result.Type)
	assert.Equal(to, result.To)
	assert.Equal(from, result.From)
	assert.Equal(summary, result.Summary)
	assert.Equal(bodyUri, result.BodyUri)
	assert.Equal(direction, result.Direction)
	assert.Equal(parts, result.Parts)
	assert.Equal(username, result.Username)
}

func TestMessagesReceived(t *testing.T) {
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

	result, err := client.Messages.Received()

	assert := assert.New(t)

	assert.Nil(err)

	assert.Equal("GET", h.Request.Method)
	assert.Equal("/v1.0/inbox/messages", h.Request.URL.String())

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

func TestMessagesReceivedWithPaging(t *testing.T) {
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

	result, err := client.Messages.Received(xesende.Page(5, 10))

	assert := assert.New(t)

	assert.Nil(err)

	assert.Equal("GET", h.Request.Method)
	assert.Equal("/v1.0/inbox/messages", h.Request.URL.Path)

	if user, pass, ok := h.Request.BasicAuth(); assert.True(ok) {
		assert.Equal("user", user)
		assert.Equal("pass", pass)
	}

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
