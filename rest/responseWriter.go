package rest

import (
	"net/http"
)

// StatusCodeResponseWriter is an implementation of http.ResponseWriter
// which can be used to obtain the response status code.
//
// Example:
//
//	func handler(w http.ResponseWriter, r *http.Request) {
//	  rw := &StatusCodeResponseWriter{}
//	  ServeHTTP( rw.Wrap( w ), r )
//	  statusCode := rw.GetStatus()
//	}
//
// The Wrap() function returns a http.ResponseWriter with supports the
// http.Flusher or http.Pusher interfaces if the one being wrapped supports it.
type StatusCodeResponseWriter struct {
	http.ResponseWriter
	code int
}

func (rw *StatusCodeResponseWriter) WriteHeader(code int) {
	rw.code = code
	rw.ResponseWriter.WriteHeader(code)
}

// GetStatus returns the staus code of the response
func (rw *StatusCodeResponseWriter) GetStatus() int {
	return rw.code
}

// Wrap wraps an existing http.ResponseWriter to our instance.
// If the wrapped instance implements the http.Flusher or http.Pusher interfaces
// then the returned instance also supports it.
func (rw *StatusCodeResponseWriter) Wrap(wrap http.ResponseWriter) http.ResponseWriter {
	rw.ResponseWriter = wrap
	flusher, _ := wrap.(http.Flusher)
	pusher, _ := wrap.(http.Pusher)

	if flusher == nil && pusher == nil {
		return rw
	}

	if flusher == nil && pusher != nil {
		return struct {
			http.ResponseWriter
			http.Pusher
		}{rw, pusher}
	}

	if flusher != nil && pusher == nil {
		return struct {
			http.ResponseWriter
			http.Flusher
		}{rw, flusher}
	}

	return struct {
		http.ResponseWriter
		http.Flusher
		http.Pusher
	}{rw, flusher, pusher}
}
