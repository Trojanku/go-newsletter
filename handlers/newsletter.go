package handlers

import (
	"Goo/model"
	"Goo/views"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
)

type signupper interface {
	SignupForNewsletter(ctx context.Context, email model.Email) (string, error)
}

type sender interface {
	Send(ctx context.Context, m model.Message) error
}

func NewsletterSignup(log *zap.Logger, mux chi.Router, s signupper, q sender) {
	mux.Post("/newsletter/signup", func(w http.ResponseWriter, r *http.Request) {

		if log == nil {
			log = zap.NewNop()
		}

		email := model.Email(r.FormValue("email"))

		if !email.IsValid() {
			http.Error(w, "email is invalid", http.StatusBadRequest)
			return
		}

		token, err := s.SignupForNewsletter(r.Context(), email)
		if err != nil {
			log.Error(fmt.Sprintf("error signing up for newsletter: %v", err))
			http.Error(w, "error signing up, refresh to try again", http.StatusBadGateway)
			return
		}

		err = q.Send(r.Context(), model.Message{
			"job":   "confirmation_email",
			"email": email.String(),
			"token": token,
		})
		if err != nil {
			log.Error(fmt.Sprintf("error sending newsletter message to queue: %v", err))
			http.Error(w, "error signing up, refresh to try again", http.StatusBadGateway)
			return
		}

		http.Redirect(w, r, "/newsletter/thanks", http.StatusFound)
	})
}

func NewsletterThanks(mux chi.Router) {
	mux.Get("/newsletter/thanks", func(w http.ResponseWriter, r *http.Request) {
		template, err := views.NewsletterThanksPage("/newsletter/thanks")
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		err = template.Execute(w, nil)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	})
}
