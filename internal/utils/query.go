package utils

import (
	"net/url"
	"strconv"
	"strings"
)

type Query struct {
	Page      int
	Limit     int
	Search    string
	SortBy    string
	SortOrder string
	Filters   map[string]string
}

func NewQuery(values url.Values) Query {
	q := Query{
		Page:      1,
		Limit:     10,
		Search:    strings.TrimSpace(values.Get("search")),
		SortBy:    strings.TrimSpace(values.Get("sortBy")),
		SortOrder: "desc",
		Filters:   make(map[string]string),
	}

	if n, err := strconv.Atoi(values.Get("page")); err == nil && n > 0 {
		q.Page = n
	}
	if n, err := strconv.Atoi(values.Get("limit")); err == nil && n > 0 && n <= 100 {
		q.Limit = n
	}

	order := strings.ToLower(strings.TrimSpace(values.Get("sortOrder")))
	if order == "asc" || order == "desc" {
		q.SortOrder = order
	}

	for key := range values {
		switch key {
		case "page", "limit", "search", "sortBy", "sortOrder":
			continue
		default:
			q.Filters[key] = strings.TrimSpace(values.Get(key))
		}
	}

	return q
}

func (q Query) Offset() int {
	return (q.Page - 1) * q.Limit
}
