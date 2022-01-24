// Package rest A basic REST server supporting HTTP.
//
// This package implements a HTTP server using net/http and github.com/gorilla/mux
// taking away most of the required boiler plate code usually needed when implementing
// basic REST services. It also provides many utility methods for handling both JSON and XML responses.
package rest

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/peter-mount/go-kernel"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"log"
	"net/http"
	"os"
	"strconv"
)

// Server The internal config of a Server
type Server struct {
	Headers       []string       // The permitted headers
	Origins       []string       // The permitted Origins
	Methods       []string       // The permitted methods
	Address       string         // Address to bind to, "" for any
	NoFlags       bool           // true to disable command line flags
	Port          int            // Port to listen to
	port          *int           // Port from command line
	router        *mux.Router    // The mux Router
	ctx           *ServerContext // Base Context
	protocol      *string        // server type
	certFile      *string        // ssl cert
	keyFile       *string        // ssl key
	logConsole    *bool          // Logging
	disableServer *bool          // Flag to disable server on command line
}

func (s *Server) Init(_ *kernel.Kernel) error {
	if !s.NoFlags {
		s.logConsole = flag.Bool("rest-log", false, "Log requests to console")
		s.protocol = flag.String("rest-protocol", "http", "Protocol to use: http|https|h2|h2c")
		s.port = flag.Int("rest-port", 0, "Port to use for http")
		s.certFile = flag.String("rest-cert", "", "TLS Certificate File")
		s.keyFile = flag.String("rest-key", "", "TLS Key File")
		s.disableServer = flag.Bool("rest-disable", false, "Disable the rest server, use for tools that run the service")
	}
	return nil
}

func (s *Server) PostInit() error {
	// Set port from command line arg or env var
	if *s.port < 1 || *s.port > 65534 {
		p, err := strconv.Atoi(os.Getenv("RESTPORT"))
		if err == nil {
			*s.port = p
		}
	}
	if *s.port > 0 && *s.port < 65535 {
		s.Port = *s.port
	}

	// Set protocol
	if *s.protocol == "" {
		*s.protocol = os.Getenv("RESTPROTOCOL")
	}
	if *s.protocol != "http" && *s.protocol != "https" && *s.protocol != "h2" && *s.protocol != "h2c" {
		return fmt.Errorf("Invalid protocol \"%s\"", *s.protocol)
	}

	if *s.certFile == "" {
		*s.certFile = os.Getenv("RESTCERT")
	}
	if *s.keyFile == "" {
		*s.keyFile = os.Getenv("RESTKEY")
	}

	s.router = mux.NewRouter()
	s.ctx = &ServerContext{context: "", server: s}

	if *s.logConsole {
		s.router.Use(ConsoleLogger())
	}

	return nil
}

// Use adds a MiddlewareHandler to the server.
// E.g. server.Use( ConsoleLogger )
func (s *Server) Use(handler mux.MiddlewareFunc) {
	s.router.Use(handler)
}

func (s *Server) Run() error {
	// Disable the server if asked
	if *s.disableServer {
		return nil
	}

	// If not defined then use port 8080
	port := s.Port
	if port < 1 || port > 65534 {
		port = 8080
	}

	// The permitted headers
	if len(s.Headers) == 0 {
		s.Headers = []string{"X-Requested-With", "Content-Type"}
	}
	if len(s.Origins) == 0 {
		s.Origins = []string{"*"}
	}
	if len(s.Methods) == 0 {
		s.Origins = []string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"}
	}
	headersOk := handlers.AllowedHeaders(s.Headers)
	originsOk := handlers.AllowedOrigins(s.Origins)
	methodsOk := handlers.AllowedMethods(s.Methods)
	handler := handlers.CORS(originsOk, headersOk, methodsOk)(s.router)

	// Now start the appropriate server
	bindingAddress := fmt.Sprintf("%s:%d", s.Address, port)
	var server *http.Server
	serveTls := false
	switch *s.protocol {
	// http/1.1
	case "http":
		serveTls = false
		server = &http.Server{
			Addr:    bindingAddress,
			Handler: handler,
		}

	// https = http/1.1 + TLS but http/2 is disabled
	case "https":
		serveTls = true
		server = &http.Server{
			Addr:    bindingAddress,
			Handler: handler,
			// This disables http/2 support
			TLSNextProto: map[string]func(*http.Server, *tls.Conn, http.Handler){},
		}

	// h2 = http/2 + TLS (also http/1.1 + TLS supported)
	case "h2":
		serveTls = true
		server = &http.Server{
			Addr:    bindingAddress,
			Handler: handler,
		}

	// h2c = http/2 with NO TLS
	//
	// See https://godoc.org/golang.org/x/net/http2/h2c#example-NewHandler
	case "h2c":
		serveTls = false
		server = &http.Server{
			Addr:    bindingAddress,
			Handler: h2c.NewHandler(handler, &http2.Server{}),
		}

	// Should not occur unless we start supporting alternate protocols
	default:
		return fmt.Errorf("Protocol %s is currently unsupported", *s.protocol)
	}

	log.Printf("Listening on %s for %s", bindingAddress, *s.protocol)
	if serveTls {
		return server.ListenAndServeTLS(*s.certFile, *s.keyFile)
	} else {
		return server.ListenAndServe()
	}
}
