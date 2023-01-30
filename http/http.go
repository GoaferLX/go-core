package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/goaferlx/go-core/log"
)

// The constants provided are defaults for net/http servers & clients.
// They are sensible default values for an API server but can and should be
// adapted by the user according to the use-case, once the object is instantiated.
// Further reading at https://blog.gopheracademy.com/advent-2016/exposing-go-on-the-internet/.
const (
	DefaultReadTimeout  time.Duration = 5 * time.Second
	DefaultWriteTimeout time.Duration = 10 * time.Second
	DefaultIdleTimeout  time.Duration = 120 * time.Second
	// DefaultClientTimeout is set the same as a server WriteTimeout so the client is not waiting for a
	// response longer than the server is going to take to send one.
	DefaultClientTimeout time.Duration = 10 * time.Second
)

// DefaultShutdown is the maximum time a process should wait for a graceful server shutdown.  Should be used for context deadlines.
const DefaultShutdownTimeout time.Duration = 10 * time.Second

type Server struct {
	*http.Server
	log.Logger
}

func (s *Server) Start(errorChan chan error) {
	s.Log("server listening", "port", s.Addr)
	errorChan <- s.ListenAndServe()
}

// Log implements the log.Logger interface.  Logging will be passed to the servers logger if one is declared, otherwise handled
// by the log package singleton.
func (s *Server) Log(msg interface{}, fields ...interface{}) error {
	if s.Logger == nil {
		return log.DefaultLogger.Log(msg, fields...)
	}
	return s.Logger.Log(msg, fields...)
}

// NewServer wraps and returns a net/http Server with pre-configured Timeouts, making it safer to
// use in production environments, as the user cannot forget to set them.
// Values used are suggested values only, the user can and should adapt them according to the use-case.
func NewServer(addr string, h http.Handler) *Server {
	return &Server{
		Server: &http.Server{
			Addr:         addr,
			Handler:      h,
			ReadTimeout:  DefaultReadTimeout,
			WriteTimeout: DefaultWriteTimeout,
			IdleTimeout:  DefaultIdleTimeout,
		},
	}

}

// NewClient returns a net/http Client with a pre-configured Timeout, making it safer to
// use in production environments, as the user cannot forget to set it.
// The value used is a suggested value only, the user can and should adapt it according to the use-case.
func NewClient() *http.Client {
	return &http.Client{
		Timeout: DefaultClientTimeout,
	}
}

// RespondWithJSON is a convenience function to set headers and write a response with a single call.
// If a response fails, the server panics.  It is good practice to wrap handlers in a recovery handler.
func RespondWithJSON(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if e, ok := data.(error); ok {
		if err := json.NewEncoder(w).Encode(e.Error()); err != nil {
			panic(err)
		}
		return
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}

}
