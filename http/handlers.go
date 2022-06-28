package http

import "net/http"

// NotFound wraps http.Error and sends a plain text error message with a http.StatusNotFound
// code.  Acts exactly like http.NotFoundHandler() but sends a custom message.
func NotFound() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, ErrNotFound.Error(), http.StatusNotFound)
	})
}

// NotAllowed wraps http.Error and sends a plain text error message with a http.MethodNotAllowed
// code.
func NotAllowed() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, ErrNotAllowed.Error(), http.StatusMethodNotAllowed)
	})
}
