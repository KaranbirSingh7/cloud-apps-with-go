package server

import "canvas/handlers"

func (s *Server) setupRoutes() {
	s.mux.HandleFunc("/", handlers.IndexPage)
	s.mux.HandleFunc("/healthz", handlers.Healthz)
}
