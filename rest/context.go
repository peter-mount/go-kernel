package rest

import (
	"github.com/gorilla/mux"
	"net/http"
)

type ServerContext struct {
	context string
	server  *Server
}

// Context creates a new ServerContext based on the current instance.
// The new Context will be the existing one suffixed with the new one.
// This allows for further grouping of rest services.
//
func (c *ServerContext) Context(context string) *ServerContext {
	if context == "" {
		return c
	}

	if context != "" && context[0:1] != "/" {
		context = "/" + context
	}

	return &ServerContext{context: c.context + context, server: c.server}
}

// Handle registers a new route with a matcher for the URL path based on the
// underlying ServerContext.
func (c *ServerContext) Handle(path string, f func(*Rest) error) *mux.Route {
	if path != "" && path[0:1] != "/" {
		path = "/" + path
	}

	return c.server.Handle(
		c.context+path,
		func(r *Rest) error {
			r.context = c.context
			return f(r)
		})
}

// HandleFunc registers a new route with a matcher for the URL path
func (c *ServerContext) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *mux.Route {
	if path != "" && path[0:1] != "/" {
		path = "/" + path
	}

	return c.server.HandleFunc(c.context+path, f)
}
