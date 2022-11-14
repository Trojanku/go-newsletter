package server

import (
	"Goo/handlers"
	"Goo/model"
	"context"
)

func (s *Server) setupRoutes() {
	handlers.Health(s.mux)
	handlers.FrontPage(s.mux)

	handlers.NewsletterSignup(s.mux, &signupperMock{})
	handlers.NewsletterThanks(s.mux)
}

type signupperMock struct{}

func (s signupperMock) SignupForNewsletter(_ context.Context, _ model.Email) (string, error) {
	return "", nil
}
