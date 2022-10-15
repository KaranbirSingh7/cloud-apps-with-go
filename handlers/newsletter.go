package handlers

import (
	"context"
	"net/http"

	"canvas/model"
	"canvas/views"

	"github.com/gorilla/mux"
)

type signupper interface {
	SignupForNewsletter(ctx context.Context, email model.Email) (string, error)
}

func NewsletterSignupWrapper(m *mux.Router, s signupper) {
	m.HandleFunc("/newsletter/signup", func(w http.ResponseWriter, r *http.Request) {
		email := model.Email(r.FormValue("email"))

		if !email.IsValid() {
			http.Error(w, "email is invalid", http.StatusBadRequest)
			return
		}

		if _, err := s.SignupForNewsletter(r.Context(), email); err != nil {
			http.Error(w, "error signing up, refresh to try again", http.StatusBadGateway)
			return
		}
		http.Redirect(w, r, "/newsletter/thanks", http.StatusFound)

	})
}

func NewsletterThanks(w http.ResponseWriter, r *http.Request) {
	_ = views.NewsletterThanksPage("/newsletter/thanks").Render(w)
}
