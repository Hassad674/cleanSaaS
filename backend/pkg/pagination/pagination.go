package pagination

import (
	"net/http"
	"strconv"
)

type Page struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Total  int `json:"total"`
}

func FromRequest(r *http.Request) Page {
	page := Page{Offset: 0, Limit: 20}

	if v := r.URL.Query().Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			page.Offset = n
		}
	}

	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 100 {
			page.Limit = n
		}
	}

	return page
}
