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

type confirmer interface {
	ConfirmNewsletterSignup(ctx context.Context, token string) (*model.Email, error)
}

func NewsletterConfirm(mux chi.Router, s confirmer, q sender) {
	mux.Get("/newsletter/confirm", func(w http.ResponseWriter, r *http.Request) {
		token := r.FormValue("token")

		template, err := views.NewsletterConfirmPage("/newsletter/confirm")
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		templateParameters := map[string]interface{}{
			"token": token,
		}
		err = template.Execute(w, templateParameters)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	})

	mux.Post("/newsletter/confirm", func(w http.ResponseWriter, r *http.Request) {
		token := r.FormValue("token")

		email, err := s.ConfirmNewsletterSignup(r.Context(), token)
		if err != nil {
			http.Error(w, "error saving email address confirmation, refresh to try again", http.StatusBadGateway)
			return
		}
		if email == nil {
			http.Error(w, "bad token", http.StatusBadRequest)
			return
		}

		err = q.Send(r.Context(), model.Message{
			"job":   "welcome_email",
			"email": email.String(),
		})
		if err != nil {
			http.Error(w, "error saving email address confirmation, refresh to try again", http.StatusBadGateway)
			return
		}
		http.Redirect(w, r, "/newsletter/confirmed", http.StatusFound)
	})
}

func NewsletterConfirmed(mux chi.Router) {
	mux.Get("/newsletter/confirmed", func(w http.ResponseWriter, r *http.Request) {
		template, err := views.NewsletterConfirmedPage("/newsletter/confirmed")
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
