package web

import (
	"compress/gzip"
	"net/http"

	"github.com/go-chi/chi/middleware"
)

func (s *Server) registerMiddlewares() {
	h := s.Mux

	h.Use(middleware.RequestID)
	h.Use(middleware.RealIP)
	h.Use(middleware.Logger)
	h.Use(middleware.Recoverer)
	h.Use(middleware.Compress(5)) //nolint:gomnd
	h.Use(gzipHandler)
}

func gzipHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("content-encoding") == "gzip" {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				panic(err)
			}
			r.Body = gz
			_ = gz.Close()
		}
		next.ServeHTTP(w, r)
	})
}
