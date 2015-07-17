package xesende

import (
	"net/http"
	"strconv"
)

type Option func(*http.Request)

func Page(startIndex, count int) Option {
	return func(r *http.Request) {
		q := r.URL.Query()

		q.Add("startindex", strconv.Itoa(startIndex))
		q.Add("count", strconv.Itoa(count))

		r.URL.RawQuery = q.Encode()
	}
}
