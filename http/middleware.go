package http

import (
	"net"
	"net/http"
	"sync"

	"github.com/goaferlx/go-core/log"
	"golang.org/x/time/rate"
)

const rateLimit = 2
const burstLimit = 4

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

// LogRequest will log the basic info of an incoming request.
func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"remote":   r.RemoteAddr,
			"protocol": r.Proto,
			"method":   r.Method,
			"URI":      r.RequestURI,
		}).Info("HTTP request")

		next.ServeHTTP(w, r)
	})
}

// RateLimit will apply a request rate limiter based on the requests IP address.
func RateLimit(next http.Handler) http.Handler {
	var mu sync.Mutex
	clients := make(map[string]*rate.Limiter)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		mu.Lock()
		limiter, ok := clients[ip]
		if !ok {
			limiter = rate.NewLimiter(rateLimit, burstLimit)
			clients[ip] = limiter
		}
		mu.Unlock()
		if !limiter.Allow() {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})

}

// RecoverPanic will attempt to recover from any panics, log the reason for the panic and return
// an internal server error.  This middleware should be applied at the start of any middleware chains.
func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.WithField("panic message", err).Error("routine paniced")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}()
		next.ServeHTTP(w, r)
	})
}
