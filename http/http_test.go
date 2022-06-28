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

		RecoverPanic(next).ServeHTTP(w, r)
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

		RecoverPanic(next).ServeHTTP(w, r)
		if got := w.Code; got != http.StatusInternalServerError {
			t.Errorf("expected status code %d, got %d", http.StatusOK, got)
		}
	})

}
