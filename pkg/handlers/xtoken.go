package handlers

import (
	"log"
	"net/http"
)

// XToken gets "X-Token" header and compare it with "t".
// Returns http.StatusForbidden if different.
func XToken(h http.Handler, t string, logger *log.Logger) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Token")
		if token != t {
			logger.Printf("[ERROR] Forbidden request: %v -> %v: wrong X-Token header: %q", r.RemoteAddr, r.URL.Path, token)
			http.Error(w, "Wrong X-Token header", http.StatusForbidden)
			return
		}

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
