package server

import (
	"Goo/messaging"
	"Goo/storage"
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	address  string
	database *storage.Database
	log      *zap.Logger
	mux      chi.Router
	queue    *messaging.Queue
	server   *http.Server
}

type Options struct {
	Database *storage.Database
	Host     string
	Log      *zap.Logger
	Port     int
	Queue    *messaging.Queue
}

func New(opts Options) *Server {
	if opts.Log == nil {
		opts.Log = zap.NewNop()
	}

	address := net.JoinHostPort(opts.Host, strconv.Itoa(opts.Port))

	// The mux is what receives an HTTP request,
	// looks at where it should go, and directs it to the code that should give a response.
	mux := chi.NewMux()

	return &Server{
		address:  address,
		database: opts.Database,
		log:      opts.Log,
		mux:      mux,
		queue:    opts.Queue,
		server: &http.Server{
			Addr:              address,
			Handler:           mux,
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      5 * time.Second,
			IdleTimeout:       5 * time.Second,
		},
	}
}

// Start the Server by setting up routes and listening for HTTP requests on the given address.
func (s *Server) Start() error {
	if err := s.database.Connect(); err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}
	s.setupRoutes()

	fmt.Println("starting on", s.address)
	s.log.Info("starting", zap.String("address", s.address))
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("error startin server: %w", err)
	}
	return nil
}

// Stop the Server gracefully within the timeout.
func (s *Server) Stop() error {
	fmt.Println("Stopping")
	s.log.Info("stopping")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("error stopping server: %w", err)
	}

	return nil
}
