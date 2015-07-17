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

func TestBatches(t *testing.T) {
	const (
		startIndex         = 0
		count              = 15
		totalCount         = 15
		id                 = "messagebatchid"
		uri                = "messagebatchuri"
		batchSize          = 1
		persistedBatchSize = 1
		accountReference   = "EXHEYEYE"
		createdBy          = "efiwewe@example.com"
		name               = "my cool batch"
	)

	var (
		status       = map[string]int{"submitted": 1}
		createdAt    = time.Date(2012, 1, 1, 12, 0, 0, 0, time.UTC)
		createdAtStr = "2012-01-01T12:00:00Z"
	)

	h := newRecordingHandler(`<?xml version="1.0" encoding="utf-8"?>
<messagebatches startindex="`+strconv.Itoa(startIndex)+`" count="`+strconv.Itoa(count)+`" totalcount="`+strconv.Itoa(totalCount)+`" xmlns="http://api.esendex.com/ns/">
 <messagebatch id="`+id+`" uri="`+uri+`">
  <createdat>`+createdAtStr+`</createdat>
  <batchsize>`+strconv.Itoa(batchSize)+`</batchsize>
  <persistedbatchsize>`+strconv.Itoa(persistedBatchSize)+`</persistedbatchsize>
  <status>
   <acknowledged>0</acknowledged>
   <authorisationfailed>0</authorisationfailed>
   <connecting>0</connecting>
   <delivered>0</delivered>
   <failed>0</failed>
   <partiallydelivered>0</partiallydelivered>
   <rejected>0</rejected>
   <scheduled>0</scheduled>
   <sent>0</sent>
   <submitted>1</submitted>
   <validityperiodexpired>0</validityperiodexpired>
   <cancelled>0</cancelled>
  </status>
  <accountreference>`+accountReference+`</accountreference>
  <createdby>`+createdBy+`</createdby>
  <name>`+name+`</name>
 </messagebatch>
</messagebatches>`, 200, map[string]string{})
	s := httptest.NewServer(h)
	defer s.Close()

	client := xesende.New("user", "pass")
	client.BaseUrl, _ = url.Parse(s.URL)

	result, err := client.Batches()

	assert := assert.New(t)

	assert.Nil(err)

	assert.Equal("GET", h.Request.Method)
	assert.Equal("/v1.1/messagebatches", h.Request.URL.String())

	if user, pass, ok := h.Request.BasicAuth(); assert.True(ok) {
		assert.Equal("user", user)
		assert.Equal("pass", pass)
	}

	assert.Equal(startIndex, result.StartIndex)
	assert.Equal(count, result.Count)
	assert.Equal(totalCount, result.TotalCount)

	if assert.Equal(1, len(result.Batches)) {
		batch := result.Batches[0]

		assert.Equal(id, batch.Id)
		assert.Equal(uri, batch.Uri)
		assert.Equal(createdAt, batch.CreatedAt)
		assert.Equal(batchSize, batch.BatchSize)
		assert.Equal(persistedBatchSize, batch.PersistedBatchSize)
		assert.Equal(status, batch.Status)
		assert.Equal(accountReference, batch.AccountReference)
		assert.Equal(createdBy, batch.CreatedBy)
		assert.Equal(name, batch.Name)
	}
}
