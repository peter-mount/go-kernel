package rest

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/peter-mount/go.uuid"
	"log"
	"net/http"
)

type IdGenerator func() (string, error)

// RequestID adds request tracing by adding the "X-Request-Id" header to the response
func RequestID(nextRequestID IdGenerator) mux.MiddlewareFunc {
	return TraceRequest("X-Request-Id", nextRequestID)
}

// TraceRequest adds a request tracking header, similar to X-Request-ID or
// X-Amz-Request-Id used by AmazonS3.
//
// header the Header name
// nextRequestID function to generate the new id
func TraceRequest(header string, nextRequestID IdGenerator) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			var err error

			// If request has the id then use it else generate one
			requestID := r.Header.Get(header)
			if requestID == "" {
				requestID, err = nextRequestID()
			}

			// No error then set it in the context & response
			if err == nil {
				ctx = context.WithValue(ctx, header, requestID)

				w.Header().Set(header, requestID)
			} else {
				log.Println("oops", err)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// DefaultIDGenerator generates a UUID
func DefaultIDGenerator() (string, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
