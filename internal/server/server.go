package server

import (
	"fmt"
	"net/http"
	"os"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(handler http.Handler) error {
	addr := fmt.Sprintf("%s:%s", os.Getenv("SERVERHOST"), os.Getenv("SERVERPORT"))
	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	return s.httpServer.ListenAndServe()
}
