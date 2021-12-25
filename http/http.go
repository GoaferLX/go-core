package http

import (
	"net/http"
	"time"
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

// NewServer returns a net/http Server with pre-configured Timeouts, making it safer to
// use in production environments, as the user cannot forget to set them.
// Values used are suggested values only, the user can and should adapt them according to the use-case.
func NewServer(addr string, h http.Handler) *http.Server {
	//logger, _ := log.ErrorLogger()
	return &http.Server{
		Addr:         addr,
		Handler:      h,
		ReadTimeout:  DefaultReadTimeout,
		WriteTimeout: DefaultWriteTimeout,
		IdleTimeout:  DefaultIdleTimeout,
		//		ErrorLog:     logger,
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

// NotFound wraps http.Error and sends a plain text error message with a http.StatusNotFound
// code.  Acts exactly like http.NotFoundHandler() but sends a custom message.
func NotFound() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, ErrNotFound.Error(), http.StatusNotFound)
	})
}
