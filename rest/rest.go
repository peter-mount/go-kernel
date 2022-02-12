package rest

import (
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
)

// A Rest query. This struct handles everything from the inbound request and
// sending the response
type Rest struct {
	writer  http.ResponseWriter
	request *http.Request
	// Response contentType
	contentType string
	// Response HTTP Status code, defaults to 200
	status int
	// The value to send
	value  interface{}
	reader io.Reader
	// Response headers
	headers map[string]string
	// true if Send() has been called
	sent bool
	// Request route variables
	vars map[string]string
	// The context
	context string
	// request attributes which are used to allow data to be stored within the
	// request whist it's being processed.
	attributes map[string]interface{}
}

// NewRest creates a new Rest query
func NewRest(writer http.ResponseWriter, request *http.Request) *Rest {
	r := &Rest{}
	r.writer = writer
	r.request = request
	r.headers = make(map[string]string)
	return r
}

// Request return the underlying http.Request so that
func (r *Rest) Request() *http.Request {
	return r.request
}

// Var returns the named route variable or "" if none
func (r *Rest) Var(name string) string {
	if r.vars == nil {
		r.vars = mux.Vars(r.request)
	}
	if r.vars == nil {
		return ""
	}
	return r.vars[name]
}

func (r *Rest) VarInt(name string, def int) int {
	s := r.Var(name)
	if s != "" {
		if i, err := strconv.Atoi(s); err == nil {
			return i
		}
	}
	return def
}

// Status sets the HTTP status of the response.
func (r *Rest) Status(status int) *Rest {
	r.status = status
	return r
}

// Value sets the response value
func (r *Rest) Value(value interface{}) *Rest {
	r.value = value
	return r
}

// Value sets the response value
func (r *Rest) Reader(rdr io.Reader) *Rest {
	r.reader = rdr
	return r
}

// Writer returns a io.Writer to write the response
func (r *Rest) Writer() io.Writer {
	// Clear any values
	r.value = nil
	r.reader = nil
	// Force a send so headers are sent
	r.Send()
	// Return the underlying writer
	return r.writer
}

// Value sets the response value
func (r *Rest) ContentType(c string) *Rest {
	r.contentType = c
	return r
}

// HTML forces the response to be html
func (r *Rest) HTML() *Rest { return r.ContentType(TEXT_HTML) }

// JSON forces the response to be JSON
func (r *Rest) JSON() *Rest { return r.ContentType(APPLICATION_JSON) }

// XML forces the response to be XML
func (r *Rest) XML() *Rest { return r.ContentType(APPLICATION_XML) }

// Context returns the base context for this request
func (r *Rest) Context() string {
	return r.context
}

// GetAttribute returns the named request attribute
func (r *Rest) GetAttribute(name string) (interface{}, bool) {
	if r.attributes == nil {
		return nil, false
	}
	v, e := r.attributes[name]
	return v, e
}

// SetAttribute returns the named request attribute
func (r *Rest) SetAttribute(n string, v interface{}) {
	if r.attributes == nil {
		r.attributes = make(map[string]interface{})
	}
	r.attributes[n] = v
}

// PushSupported returns true if http2 Push is supported
func (r *Rest) PushSupported() bool {
	_, ok := r.writer.(http.Pusher)
	return ok
}

// Push initiates an HTTP/2 server push. This constructs a synthetic
// request using the given target and options, serializes that request
// into a PUSH_PROMISE frame, then dispatches that request using the
// server's request handler. If opts is nil, default options are used.
//
// The target must either be an absolute path (like "/path") or an absolute
// URL that contains a valid host and the same scheme as the parent request.
// If the target is a path, it will inherit the scheme and host of the
// parent request.
//
// The HTTP/2 spec disallows recursive pushes and cross-authority pushes.
// Push may or may not detect these invalid pushes; however, invalid
// pushes will be detected and canceled by conforming clients.
//
// Handlers that wish to push URL X should call Push before sending any
// data that may trigger a request for URL X. This avoids a race where the
// client issues requests for X before receiving the PUSH_PROMISE for X.
//
// Push returns ErrNotSupported if the client has disabled push or if push
// is not supported on the underlying connection.
func (r *Rest) Push(target string, opts *http.PushOptions) error {
	p, ok := r.writer.(http.Pusher)
	if ok {
		return p.Push(target, opts)
	}
	return nil
}
