// Package server contains everything for setting up and running the HTTP server.
package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"go.uber.org/zap"
)

type Server struct {
	address string
	server  *http.Server
	log     *zap.Logger
	mux     *mux.Router
}

type Options struct {
	Host string
	Port int
	Log  *zap.Logger
}

func (o *Options) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		o.Log.Info("", zap.String("url", r.RequestURI))
		next.ServeHTTP(w, r)
	})
}

func NewServer(opts Options) *Server {
	if opts.Log == nil {
		opts.Log = zap.NewNop()
	}

	address := net.JoinHostPort(opts.Host, strconv.Itoa(opts.Port))

	r := mux.NewRouter()

	r.Use(opts.loggingMiddleware)

	return &Server{
		address: address,
		mux:     r,
		server: &http.Server{
			Addr:              net.JoinHostPort("", strconv.Itoa(opts.Port)),
			Handler:           r,
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      5 * time.Second,
			IdleTimeout:       5 * time.Second,
		},
		log: opts.Log,
	}

}

func (s *Server) Start() error {
	s.setupRoutes()

	// log.Printf("Starting on %s", s.address)
	s.log.Info("Starting", zap.String("address", s.address))
	if err := s.server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("error starting server: %w", err)
		}
	}
	return nil
}

func (s *Server) Stop() error {
	s.log.Info("Stopping")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("error stopping server: %w", err)
	}

	return nil
}
