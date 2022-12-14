package server

import (
	"Goo/handlers"
)

func (s *Server) setupRoutes() {
	handlers.Health(s.mux, s.database)
	handlers.FrontPage(s.mux)

	handlers.NewsletterSignup(s.log, s.mux, s.database, s.queue)
	handlers.NewsletterThanks(s.mux)

	handlers.NewsletterConfirm(s.mux, s.database, s.queue)
	handlers.NewsletterConfirmed(s.mux)
}
