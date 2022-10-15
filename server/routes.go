package server

import (
	"canvas/handlers"
)

func (s *Server) setupRoutes() {
	s.mux.HandleFunc("/", handlers.IndexPage)
	s.mux.HandleFunc("/healthz", handlers.Healthz)
	s.mux.HandleFunc("/newsletter/signup", handlers.NewsletterSignup)
	s.mux.HandleFunc("/newsletter/thanks", handlers.NewsletterThanks)
}
