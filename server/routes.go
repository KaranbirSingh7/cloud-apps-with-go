package server

import (
	"canvas/handlers"
	"canvas/model"
	"context"
)

func (s *Server) setupRoutes() {
	s.mux.HandleFunc("/", handlers.IndexPage)

	handlers.HealthzWrapper(s.mux, s.database)
	// s.mux.HandleFunc("/newsletter/signup", handlers.NewsletterSignup)
	handlers.NewsletterSignupWrapper(s.mux, &signupperMock{}) //handle: /newsletter/signup URL
	s.mux.HandleFunc("/newsletter/thanks", handlers.NewsletterThanks)

}

type signupperMock struct{}

func (s signupperMock) SignupForNewsletter(ctx context.Context, email model.Email) (string, error) {
	return "", nil
}
