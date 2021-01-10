package rest

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// HandleFunc registers a new route with a matcher for the URL path
func (s *Server) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *mux.Route {
	if path != "" && path[0:1] != "/" {
		path = "/" + path
	}

	return s.router.HandleFunc(path, f)
}

// Handle registers a new route with a matcher for the URL path.
// Unlike HandleFunc() this accepts a func that just takes a Rest instance.
// If the returned error is not nil then a 500 response is issued otherwise
// the response is sent unless Rest.Send() has already been called.
func (s *Server) Handle(path string, f func(*Rest) error) *mux.Route {
	if path != "" && path[0:1] != "/" {
		path = "/" + path
	}

	return s.HandleFunc(path, Handler(f))
}

// Handler creates a wrapper around a rest handler and one used by mux
func Handler(f func(*Rest) error) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		rest := NewRest(w, r)

		if err := f(rest); err != nil {
			// Respond with 500 no body
			log.Println(err)
			w.WriteHeader(500)
		} else {
			// Send the response
			err := rest.Send()
			if err != nil {
				log.Println(err)
			}
		}
	}
}

// NotFound Adds a custom NotFound handler
func (s *Server) NotFound(f func(*Rest) error) {
	s.router.NotFoundHandler = http.HandlerFunc(Handler(f))
}
