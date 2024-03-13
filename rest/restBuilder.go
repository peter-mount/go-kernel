package rest

import (
	"github.com/gorilla/mux"
	"net/http"
)

type RestHandler func(*Rest) error

type RestDecorator func(RestHandler) RestHandler

type RestBuilder struct {
	server          *Server
	commonDecorator []RestDecorator
	// pathPrefix
	pathPrefix      string
	pathPrefixRoute *mux.Router
	// Below this point are reset after every Build()
	authDecorator []RestDecorator
	handler       RestHandler
	headers       []string
	headersRegexp []string
	methods       []string
	paths         []string
	queries       []string
}

func (s *Server) RestBuilder() *RestBuilder {
	return &RestBuilder{server: s}
}

// Build a rest endpoint
func (r *RestBuilder) Build() *RestBuilder {
	// Form the final handler with any decorators applied
	h := r.handler

	// Apply auth decorators
	for _, d := range r.authDecorator {
		h = d(h)
	}

	// Apply common decorators
	for _, d := range r.commonDecorator {
		h = d(h)
	}

	// The final handler
	hf := http.HandlerFunc(Handler(h))

	if len(r.paths) == 0 {
		// Rare but could happen - e.g. we are not matching against a path but a Header or Query string
		r.buildRoute(r.newRoute().Handler(hf))
	} else {
		for _, path := range r.paths {
			r.buildRoute(r.newRoute().Handler(hf).Path(r.pathPrefix + path))
		}
	}

	return r.Reset()
}

func (r *RestBuilder) newRoute() *mux.Route {
	if r.pathPrefix != "" {
		return r.pathPrefixRoute.NewRoute()
	}
	return r.server.router.NewRoute()
}

func (r *RestBuilder) buildRoute(route *mux.Route) {
	r.applyOption(r.methods, route.Methods)
	r.applyOption(r.headers, route.Headers)
	r.applyOption(r.headersRegexp, route.HeadersRegexp)
	r.applyOption(r.queries, route.Queries)
}

func (r *RestBuilder) applyOption(a []string, f func(...string) *mux.Route) {
	if len(a) > 0 {
		f(a...)
	}
}

// Reset the builder so only the common entries are used in the next
// Build(), e.g. Decorator()'s' and PathPrefix() are affected by this.
//
// Note, Build() call this so it's not normally required
func (r *RestBuilder) Reset() *RestBuilder {
	// Reset the builder
	r.authDecorator = []RestDecorator{}
	r.handler = nil
	r.headers = []string{}
	r.headersRegexp = []string{}
	r.methods = []string{}
	r.paths = []string{}
	r.queries = []string{}

	return r
}

// PathPrefix is the common prefix that will be prepended to all endpoints
// in this builder.
//
// It's preferable if a group of endpoints share the same
// contant prefix, e.g. "/api" to use this not just to save typing the same
// prefix as this will improve performance during matching.
//
// Any trailing "/" are removed.
//
// If the prefix does not start with a "/" then one is prepended.
//
// "" or "/" are treated as no prefix.
func (r *RestBuilder) PathPrefix(pathPrefix string) *RestBuilder {
	if pathPrefix == "/" {
		pathPrefix = ""
	}

	if pathPrefix != "" {
		if pathPrefix[0:1] != "/" {
			pathPrefix = "/" + pathPrefix
		}
		for len(pathPrefix) > 1 && pathPrefix[len(pathPrefix)-1:] == "/" {
			pathPrefix = pathPrefix[0 : len(pathPrefix)-1]
		}
	}

	if pathPrefix == "/" {
		pathPrefix = ""
	}

	r.pathPrefix = pathPrefix

	// Create the router for this prefix
	if pathPrefix != "" {
		r.pathPrefixRoute = r.server.router.NewRoute().
			PathPrefix(pathPrefix + "/").
			Subrouter()
	}

	return r
}

// Path to be appended to the end point being built.
// You can provide multiple patterns for the call being built
func (r *RestBuilder) Path(s ...string) *RestBuilder {
	r.paths = s

	// Ensure paths start with "/"
	for i, path := range r.paths {
		if path != "" && path[0:1] != "/" {
			r.paths[i] = "/" + path
		}
	}

	return r
}

// Headers adds a matcher for request header values.
// It accepts a sequence of key/value pairs to be matched. For example:
//
//	r := server.RestBuilder()
//	r.Headers("Content-Type", "application/json",
//	          "X-Requested-With", "XMLHttpRequest")
//
// The above route will only match if both request header values match.
// If the value is an empty string, it will match any value if the key is set.
func (r *RestBuilder) Headers(s ...string) *RestBuilder {
	r.headers = s
	return r
}

// HeadersRegexp accepts a sequence of key/value pairs, where the value has regex
// support. For example:
//
//	r := server.RestBuilder()
//	r.HeadersRegexp("Content-Type", "application/(text|json)",
//	          "X-Requested-With", "XMLHttpRequest")
//
// The above route will only match if both the request header matches both regular expressions.
// If the value is an empty string, it will match any value if the key is set.
// Use the start and end of string anchors (^ and $) to match an exact value.
func (r *RestBuilder) HeadersRegexp(s ...string) *RestBuilder {
	r.headersRegexp = s
	return r
}

// Method sets the HTTP Method this endpoint is to use
func (r *RestBuilder) Method(s ...string) *RestBuilder {
	r.methods = s
	return r
}

// Queries adds a matcher for URL query values.
// It accepts a sequence of key/value pairs. Values may define variables.
// For example:
//
//	r := server.RestBuilder()
//	r.Queries("foo", "bar", "id", "{id:[0-9]+}")
//
// The above route will only match if the URL contains the defined queries
// values, e.g.: ?foo=bar&id=42.
//
// It the value is an empty string, it will match any value if the key is set.
//
// Variables can define an optional regexp pattern to be matched:
//
// - {name} matches anything until the next slash.
//
// - {name:pattern} matches the given regexp pattern.
func (r *RestBuilder) Queries(s ...string) *RestBuilder {
	r.queries = s
	return r
}

// Handler sets the RestHandler to use for this endpoint
func (r *RestBuilder) Handler(f RestHandler) *RestBuilder {
	r.handler = f
	return r
}

// Decorate will invoke a RestDecorator against the RestHandler set by Handler()
// to decorate it. This is usually used to add error handling,
// ensuring certain headers are in the response etc.
// Note: This decorator will be used by all endpoints built by this instance
// of the builder. To apply a decorator on a single endpoint use Authenticator()
func (r *RestBuilder) Decorate(f RestDecorator) *RestBuilder {
	r.commonDecorator = append(r.commonDecorator, f)
	return r
}

// Authenticator is a RestDecorator which is specific to the current endpoint
// being built. It's called Authenticator as it's usually used for authentication
// but it can be used for other purposes.
func (r *RestBuilder) Authenticator(f RestDecorator) *RestBuilder {
	r.authDecorator = append(r.authDecorator, f)
	return r
}

// AddHeadersDecorator is a decorator which ensures that the given headers
// are always present in a response that has not returned an error.
//
// Here's an example on ensuring certain headers are always present:
//
//	Decorate( (&rest.AddHeadersDecorator{
//	  "Access-Control-Allow-Origin": "*",
//	  "X-Clacks-Overhead": "GNU Terry Pratchett",
//	}).Decorator )
type AddHeadersDecorator map[string]string

func (a *AddHeadersDecorator) Decorator(h RestHandler) RestHandler {
	return func(r *Rest) error {
		err := h(r)
		if err == nil {
			for k, v := range *a {
				r.AddHeader(k, v)
			}
		}
		return err
	}
}
