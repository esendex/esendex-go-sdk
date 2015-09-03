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

func TestAccounts(t *testing.T) {
	const (
		id                = "accountid"
		uri               = "http://someaccount"
		reference         = "EX093052"
		label             = "My account"
		address           = "443523"
		accountType       = "Professional"
		messagesRemaining = 2322
		role              = "CoolUser"
		settingsUri       = "http://somesettings"
	)

	var (
		expiresOn    = time.Date(2012, 1, 1, 12, 0, 5, 0, time.UTC)
		expiresOnStr = "2012-01-01T12:00:05"
	)

	h := newRecordingHandler(`<?xml version="1.0" encoding="utf-8"?>
<accounts xmlns="http://api.esendex.com/ns/">
 <account id="`+id+`" uri="`+uri+`">
  <reference>`+reference+`</reference>
  <label>`+label+`</label>
  <address>`+address+`</address>
  <type>`+accountType+`</type>
  <messagesremaining>`+strconv.Itoa(messagesRemaining)+`</messagesremaining>
  <expireson>`+expiresOnStr+`</expireson>
  <role>`+role+`</role>
  <settings uri="`+settingsUri+`" />
 </account>
</accounts>`, 200, map[string]string{})
	s := httptest.NewServer(h)
	defer s.Close()

	client := xesende.New("user", "pass")
	client.BaseURL, _ = url.Parse(s.URL)

	result, err := client.Accounts()

	assert := assert.New(t)

	assert.Nil(err)

	assert.Equal("GET", h.Request.Method)
	assert.Equal("/v1.0/accounts", h.Request.URL.String())

	if user, pass, ok := h.Request.BasicAuth(); assert.True(ok) {
		assert.Equal("user", user)
		assert.Equal("pass", pass)
	}

	if assert.Equal(1, len(result.Accounts)) {
		account := result.Accounts[0]

		assert.Equal(id, account.ID)
		assert.Equal(uri, account.URI)
		assert.Equal(reference, account.Reference)
		assert.Equal(label, account.Label)
		assert.Equal(address, account.Address)
		assert.Equal(accountType, account.Type)
		assert.Equal(messagesRemaining, account.MessagesRemaining)
		assert.Equal(expiresOn, account.ExpiresOn)
		assert.Equal(role, account.Role)
		assert.Equal(settingsUri, account.SettingsURI)
	}
}
