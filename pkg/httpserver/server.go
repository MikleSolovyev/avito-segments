package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

func New(handler http.Handler, cfg *Config) *Server {
	httpServer := &http.Server{
		Handler:      handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		Addr:         cfg.Address,
	}

	s := &Server{
		server:          httpServer,
		notify:          make(chan error, 1),
		shutdownTimeout: cfg.ShutdownTimeout,
	}

	s.start()

	return s
}

func (s *Server) start() {
	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	const fn = "httpserver.Shutdown"

	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	err := s.server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}
