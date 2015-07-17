package xesende

import (
	"net/http"
	"strconv"
)

// Option is a function that mutates a request.
type Option func(*http.Request)

// Page is an option that sets the startindex and count query parameters.
func Page(startIndex, count int) Option {
	return func(r *http.Request) {
		q := r.URL.Query()

		q.Add("startindex", strconv.Itoa(startIndex))
		q.Add("count", strconv.Itoa(count))

		r.URL.RawQuery = q.Encode()
	}
}
