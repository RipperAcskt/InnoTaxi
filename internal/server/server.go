package server

import "net/http"

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	return s.httpServer.ListenAndServe()
}
