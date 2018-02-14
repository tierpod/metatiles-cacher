package main

import (
	"log"
	"net/http"
)

type tokenStore struct {
	token  string
	logger *log.Logger
}

func newTokenStore(token string, logger *log.Logger) *tokenStore {
	return &tokenStore{
		token:  token,
		logger: logger,
	}
}

func (t *tokenStore) Middleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Token")
		if token != t.token {
			t.logger.Printf("[ERROR] (xtoken) forbidden request from %v to %v: wrong X-Token header: %q", r.RemoteAddr, r.URL.Path, token)
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
