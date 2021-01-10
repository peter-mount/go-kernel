package rest

import (
	"github.com/gorilla/mux"
	"net/http"
)

// Get returns a named route.
// You can name a Route by calling Name(name string) against the route returned
// by either Handle or HandleFunc
func (s *Server) Get(name string) *mux.Route {
	return s.router.Get(name)
}

// Get returns a named route.
// You can name a Route by calling Name(name string) against the route returned
// by either Handle or HandleFunc
func (c *ServerContext) Get(name string) *mux.Route {
	return c.server.Get(name)
}

// Static adds a static file service at a specific prefix.
// prefix is the prefix to serve from, e.g. "/static/"
// dir is the directory to serve static content from
func (s *Server) Static(prefix, dir string) {
	fileServer := http.FileServer(http.Dir(dir))
	s.router.PathPrefix(prefix).
		Handler(http.StripPrefix(prefix, fileServer))
}
