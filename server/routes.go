package server

import (
	"Goo/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) setupRoutes() {
	s.mux.Use(handlers.AddMetrics(s.metrics))

	handlers.Health(s.mux, s.database)

	handlers.FrontPage(s.mux)
	handlers.NewsletterSignup(s.mux, s.database, s.queue, s.log)
	handlers.NewsletterThanks(s.mux)
	handlers.NewsletterConfirm(s.mux, s.database, s.queue, s.log)
	handlers.NewsletterConfirmed(s.mux)

	s.mux.Group(func(r chi.Router) {
		r.Use(middleware.BasicAuth("goo", map[string]string{"admin": s.adminPassword}))

		handlers.MigrateTo(r, s.database)
		handlers.MigrateUp(r, s.database)
	})

	metricsAuth := middleware.BasicAuth("metrics", map[string]string{"prometheus": s.metricsPassword})
	handlers.Metrics(s.mux.With(metricsAuth), s.metrics)
}
