package server

import (
	"net/http"

	"github.com/RipperAcskt/innotaxi/config"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(handler http.Handler, cfg *config.Config) error {
	s.httpServer = &http.Server{
		Addr:    cfg.SERVER_HOST,
		Handler: handler,
	}

	return s.httpServer.ListenAndServe()
}
