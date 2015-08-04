package xesende

import "net/http"

func (c *AccountClient) Batches(opts ...Option) (*BatchesResponse, error) {
	accountOption := func(r *http.Request) {
		q := r.URL.Query()

		q.Add("filterBy", "account")
		q.Add("filterValue", c.reference)

		r.URL.RawQuery = q.Encode()
	}

	return c.Client.Batches(append(opts, accountOption)...)
}
