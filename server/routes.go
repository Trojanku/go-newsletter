package server

import "Goo/handlers"

func (s *Server) setupRoutes() {
	handlers.Health(s.mux)
}