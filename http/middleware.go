package http

import (
	"net/http"

	"github.com/goaferlx/go-core/log"
)

// CheckContentHeader checks the incoming request Content-Type Header matches a user specified header before handling the request.
// If false it sets a 415 MediaNotSupported return header and processing is stopped.
// An empty string is not considered to be false.  GET requests are ignored as no content is sent.
func CheckContentHeader(contentType string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				next.ServeHTTP(w, r)
				return
			}
			if header := r.Header.Get("Content-Type"); header != "" && header != contentType {
				w.WriteHeader(http.StatusUnsupportedMediaType)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// CheckAcceptHeader checks the incoming request Accept header matches a user specified header before handling the request.
// If false it writes 406 Not Acceptable return header and processing is stopped.
// If true, or the client does not send a header or will accept all content, execution continues.
func CheckAcceptHeader(acceptType string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if header := r.Header.Get("Accept"); header != "*/*" && header != acceptType {
				w.WriteHeader(http.StatusNotAcceptable)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// Do I want this in this package?
func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.HTTPRequest(r)
		next.ServeHTTP(w, r)
	})
}
