package http

import "net/http"

// CheckContentType checks the incoming request Content-Type Header matches a user specified header before handling the request.
// If false it sets a return header and processing is stopped.
// An empty string is not considered to be false.  GET requests are ignored as no content is sent.
func CheckContentType(contentType string) func(next http.Handler) http.Handler {
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
