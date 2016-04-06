package gitest

import (
	"net/http"

	"github.com/heroku/pat"
)

// Handler is all the http routes required
func (s *Server) Handler() *pat.PatternServeMux {
	h := pat.New()
	h.GetFunc("/:name.git/info/refs", s.accessMiddleware(s.refsEndpoint))
	h.PostFunc("/:name.git/:service", s.accessMiddleware(s.serviceEndpoint))

	return h
}

func (s *Server) accessMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var name = r.URL.Query().Get(":name")

		if name == s.ValidRepo {
			next.ServeHTTP(w, r)
		} else if name == s.NotAllowedRepo {
			w.WriteHeader(401)
		} else {
			w.WriteHeader(404)
		}
	})
}
