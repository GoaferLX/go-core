package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecoverPanic(t *testing.T) {
	t.Run("does nothing if no panic", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/irrelevant", nil)

		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
		})
		srv := NewServer("", nil)
		srv.RecoverPanic(next).ServeHTTP(w, r)
		if got := w.Code; got != http.StatusCreated {
			t.Errorf("expected status code %d, got %d", http.StatusOK, got)
		}
	})

	t.Run("returns 500 if panic happens", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/irrelevant", nil)

		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("something happened")
		})
		srv := NewServer("", nil)
		srv.RecoverPanic(next).ServeHTTP(w, r)
		if got := w.Code; got != http.StatusInternalServerError {
			t.Errorf("expected status code %d, got %d", http.StatusOK, got)
		}
	})

}

type mockLimiter struct {
	maxRequests     int
	currentRequests int
}

func (ml *mockLimiter) Allow() bool {
	if ml.currentRequests < ml.maxRequests {
		ml.currentRequests++
		return true
	}
	return false
}
func TestRateLimit(t *testing.T) {
	t.Run("new request generates a new client", func(t *testing.T) {

		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(404)
		})
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/irrelevant", nil)

		clients := make(map[string]Allower)

		RateLimit(clients)(next).ServeHTTP(w, r)

		if len(clients) != 1 {
			t.Errorf("expected 1 item in clients, got %d", len(clients))
		}

	})
	t.Run("requests from different ips generate different clients", func(t *testing.T) {

		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
		})
		clients := make(map[string]Allower)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/irrelevant", nil)

		r.RemoteAddr = "127.0.0.1:3000"
		RateLimit(clients)(next).ServeHTTP(w, r)

		r.RemoteAddr = "127.0.0.2:3000"
		RateLimit(clients)(next).ServeHTTP(w, r)

		if len(clients) != 2 {
			t.Errorf("expected 2 items in clients, got %d", len(clients))
		}

	})

	t.Run("request from known source uses exisiting client", func(t *testing.T) {
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(404)
		})
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/irrelevant", nil)

		r.RemoteAddr = "127.0.0.1:3000"
		clients := make(map[string]Allower)
		limiter := &mockLimiter{currentRequests: 0, maxRequests: 1}
		clients["127.0.0.1"] = limiter

		RateLimit(clients)(next).ServeHTTP(w, r)

		if len(clients) != 1 {
			t.Errorf("expected 1 item in clients, got %d", len(clients))
		}
		if limiter.currentRequests != 1 {
			t.Errorf("limiter has not been called")
		}

	})

	t.Run("returns 429 if max requests exceeded", func(t *testing.T) {
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(404)
		})
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/irrelevant", nil)

		r.RemoteAddr = "127.0.0.1:3000"
		clients := make(map[string]Allower)
		clients["127.0.0.1"] = &mockLimiter{currentRequests: 2, maxRequests: 2}

		RateLimit(clients)(next).ServeHTTP(w, r)

		if len(clients) != 1 {
			t.Errorf("expected 1 item in clients, got %d", len(clients))
		}
		expected := 429
		if got := w.Code; got != expected {
			t.Errorf("expected status code %d, got %d", expected, got)
		}

	})
}
