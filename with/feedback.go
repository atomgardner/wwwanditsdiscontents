package with

import (
	"log"
	"net/http"
)

func Feedback(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// XXX: could make this more generic with a strings.Replacer{}
		log.Printf("[%s] %s %s %s", r.RemoteAddr, r.Host, r.RequestURI, r.URL.Path)
		h.ServeHTTP(w, r)
	})
}
