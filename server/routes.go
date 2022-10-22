package server

import (
	"canvas/handlers"
)

func (s *Server) setupRoutes() {
	s.mux.HandleFunc("/", handlers.IndexPage)
	handlers.HealthzWrapper(s.mux, s.database)
	// s.mux.HandleFunc("/newsletter/signup", handlers.NewsletterSignup)
	handlers.NewsletterSignupWrapper(s.mux, s.database, s.queue) //handle: /newsletter/signup URL
	s.mux.HandleFunc("/newsletter/thanks", handlers.NewsletterThanks)

}
