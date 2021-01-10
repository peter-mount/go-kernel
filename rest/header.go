package rest

import (
	"fmt"
)

// AddHeader adds a header to the response, replacing any existing entry
func (r *Rest) AddHeader(name string, value string) *Rest {
	r.headers[name] = value
	return r
}

// GetHeader returns a header from the request or "" if not present
func (r *Rest) GetHeader(name string) string {
	return r.request.Header.Get(name)
}

func (r *Rest) AccessControlAllowOrigin(value string) *Rest {
	if value == "" {
		value = "*"
	}
	return r.AddHeader("Access-Control-Allow-Origin", value)
}

func (r *Rest) Etag(tag string) *Rest {
	return r.AddHeader("Etag", "\""+tag+"\"")
}

func (r *Rest) CacheMaxAge(age int) *Rest {
	return r.AddHeader("Cache-Control", fmt.Sprintf("public, max-age=%d, s-maxage=%d", age, age))
}

func (r *Rest) CacheNoCache() *Rest {
	return r.AddHeader("Cache-Control", "no-cache, no-store, must-revalidate")
}
