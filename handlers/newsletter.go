package handlers

import (
	"context"
	"net/http"

	"canvas/model"
	"canvas/views"
)

type signupper interface {
	SignupForNewsletter(ctx context.Context, email model.Email) (string, error)
}

type signupperMock struct{}

func (s signupperMock) SignupForNewsletter(ctx context.Context, email model.Email) (string, error) {
	return "", nil
}

func NewsletterSignup(w http.ResponseWriter, r *http.Request) {
	email := model.Email(r.FormValue("email"))

	var s signupperMock
	if !email.IsValid() {
		http.Error(w, "email is invalid", http.StatusBadRequest)
		return
	}

	if _, err := s.SignupForNewsletter(r.Context(), email); err != nil {
		http.Error(w, "error signing up, refresh to try again", http.StatusBadGateway)
		return
	}
	http.Redirect(w, r, "/newsletter/thanks", http.StatusFound)
}

func NewsletterThanks(w http.ResponseWriter, r *http.Request) {
	_ = views.NewsletterThanksPage("/newsletter/thanks").Render(w)
}
