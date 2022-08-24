package http

import (
	"net"
	"net/http"
	"sync"

	"github.com/goaferlx/go-core/log"
	"golang.org/x/time/rate"
)

const DefaultRateLimit = 2
const DefaultBurstLimit = 4

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
			"remote_address": r.RemoteAddr,
			"protocol":       r.Proto,
			"method":         r.Method,
			"URI":            r.RequestURI,
		}).Info("HTTP request")

		next.ServeHTTP(w, r)
	})
}

type Allower interface {
	Allow() bool
}

// RateLimit will apply a request rate limiter based on the incoming requests IP address.
func RateLimit(clients map[string]Allower) func(next http.Handler) http.Handler {
	mw := RateLimiterMw{
		Limit:   DefaultRateLimit,
		Burst:   DefaultBurstLimit,
		clients: clients,
	}
	return mw.RateLimit()
}

type RateLimiterMw struct {
	Limit   int
	Burst   int
	clients map[string]Allower
	mu      sync.Mutex
}

func (mw *RateLimiterMw) RateLimit() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			mw.mu.Lock()
			limiter, ok := mw.clients[ip]
			if !ok {
				limiter = rate.NewLimiter(rate.Limit(mw.Limit), mw.Burst)
				mw.clients[ip] = limiter
			}
			mw.mu.Unlock()

			if !limiter.Allow() {
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
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
