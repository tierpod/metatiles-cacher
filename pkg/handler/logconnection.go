package handler

import (
	"log"
	"net/http"
)

// LogConnection logs information about http connection
// TODO: do something like this?
// 10.145.0.45 - - [19/Oct/2017:03:08:01 +0500] "GET /style/16/60504/3936.png HTTP/1.1" 200 126
func LogConnection(h http.Handler, l *log.Logger) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		l.Printf("%v - \"%v %v\"", r.RemoteAddr, r.Method, r.URL.Path)
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
