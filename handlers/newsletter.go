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

type sender interface {
	Send(ctx context.Context, m model.Message) error
}

func NewsletterSignupWrapper(m *mux.Router, s signupper, q sender) {
	m.HandleFunc("/newsletter/signup", func(w http.ResponseWriter, r *http.Request) {
		email := model.Email(r.FormValue("email"))

		if !email.IsValid() {
			http.Error(w, "email is invalid", http.StatusBadRequest)
			return
		}

		token, err := s.SignupForNewsletter(r.Context(), email)
		if err != nil {
			http.Error(w, "error signing up, refresh to try again", http.StatusBadGateway)
			return
		}

		err = q.Send(r.Context(), model.Message{
			"job":   "confirmation_email",
			"email": email.String(),
			"token": token,
		})
		if err != nil {
			http.Error(w, "error signing up, refresh to try again", http.StatusBadGateway)
			return
		}

		http.Redirect(w, r, "/newsletter/thanks", http.StatusFound)

	})
}

func NewsletterThanks(w http.ResponseWriter, r *http.Request) {
	_ = views.NewsletterThanksPage("/newsletter/thanks").Render(w)
}
