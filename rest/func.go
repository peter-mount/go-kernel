package rest

import (
	"context"
	"github.com/gorilla/mux"
)

const (
	restKey = "kernel.rest.key"
)

// Do uses a func(Context) style handler for the given path.
func (s *Server) Do(path string, f func(ctx context.Context) error) *mux.Route {
	return s.Handle(path, func(rest *Rest) error {
		return f(context.WithValue(rest.Request().Context(), restKey, rest))
	})
}

// GetRest returns the *Rest instance from a Context issued with the Do() function.
func GetRest(ctx context.Context) *Rest {
	return ctx.Value(restKey).(*Rest)
}
