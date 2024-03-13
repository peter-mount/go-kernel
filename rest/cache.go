// Package rest is a basic REST server supporting HTTP.
//
// This package implements a HTTP server using net/http and github.com/gorilla/mux
// taking away most of the required boilerplate code usually needed when implementing
// basic REST services. It also provides many utility methods for handling both JSON and XML responses.
package rest

import (
	"fmt"
)

// CacheControl manages the Cache-Control header in the response.
// cache < 0 for "no cache",
// cache = 0 to ignore (i.e. don't set the header)
// cache > 0 Time in seconds for the client to cache the response
func (r *Rest) CacheControl(cache int) *Rest {
	if cache < 0 {
		r.AddHeader("Cache-Control", "no-cache")
	} else if cache > 0 {
		r.AddHeader("Cache-Control", fmt.Sprintf("max-age=%d, s-maxage=%d", cache, cache))
	}
	return r
}
