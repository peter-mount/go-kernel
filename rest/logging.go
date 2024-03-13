package rest

import (
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	// The log format:
	// | Status | Method | Host | RemoteAddr | Path | ElapsedTime | X-Request-Id
	logFormat = "| %d | %s | %s | %v | %v | %.6fs | %s"
)

// ConsoleLogger returns a middleware function to log requests to the console
func ConsoleLogger() mux.MiddlewareFunc {
	return FormatLogger(log.Printf)
}

// LogLogger returns a middleware function to log requests to a specific logger
func LogLogger(l *log.Logger) mux.MiddlewareFunc {
	return FormatLogger(l.Printf)
}

type FormatLoggerFunc func(string, ...interface{})

// FormatLogger returns a middleware function to log requests to a function.
// Note: The format passed to this function will not have a trailing "\n"
// nor will it have any time information in it.
func FormatLogger(f FormatLoggerFunc) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rw := &StatusCodeResponseWriter{}

			defer func() {
				elapsed := time.Now().Sub(start).Seconds()

				f(logFormat, rw.GetStatus(), r.Method, r.Host, remoteAddr(r), r.URL, elapsed, w.Header().Get("X-Request-Id"))
			}()

			next.ServeHTTP(rw.Wrap(w), r)
		})
	}
}

// Get the remote address, accounting for intermediate proxies
func remoteAddr(r *http.Request) net.IP {
	forwardedFor := r.Header.Get("X-Forwarded-For")
	if forwardedFor != "" {
		// Syntax on MDN: X-Forwarded-For: <client>, <proxy1>, <proxy2>
		ip := strings.Split(forwardedFor, ", ")[0]
		if !(strings.HasPrefix(ip, "10.") ||
			strings.HasPrefix(ip, "192.168.") ||
			strings.HasPrefix(ip, "172.16.") ||
			strings.HasPrefix(ip, "172.17.") ||
			strings.HasPrefix(ip, "172.18.") ||
			strings.HasPrefix(ip, "172.19.") ||
			strings.HasPrefix(ip, "172.20.") ||
			strings.HasPrefix(ip, "172.21.") ||
			strings.HasPrefix(ip, "172.22.") ||
			strings.HasPrefix(ip, "172.23.") ||
			strings.HasPrefix(ip, "172.24.") ||
			strings.HasPrefix(ip, "172.25.") ||
			strings.HasPrefix(ip, "172.26.") ||
			strings.HasPrefix(ip, "172.27.") ||
			strings.HasPrefix(ip, "172.28.") ||
			strings.HasPrefix(ip, "172.29.") ||
			strings.HasPrefix(ip, "172.30.") ||
			strings.HasPrefix(ip, "172.31.")) {
			return net.ParseIP(ip)
		}
	}

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return net.ParseIP(ip)
}
