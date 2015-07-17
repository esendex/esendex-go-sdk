package xesende

import (
	"net/http"
	"strconv"
	"time"
)

// Option is a function that mutates a request.
type Option func(*http.Request)

// Page creates an option that sets the startindex and count query parameters.
func Page(startIndex, count int) Option {
	return func(r *http.Request) {
		q := r.URL.Query()

		q.Add("startindex", strconv.Itoa(startIndex))
		q.Add("count", strconv.Itoa(count))

		r.URL.RawQuery = q.Encode()
	}
}

// Between creates an option that sets the start and finish query parameters.
func Between(start, finish time.Time) Option {
	return func(r *http.Request) {
		q := r.URL.Query()

		q.Add("start", start.Format(time.RFC3339))
		q.Add("finish", finish.Format(time.RFC3339))

		r.URL.RawQuery = q.Encode()
	}
}
