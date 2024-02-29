package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/kiennyo/syncwatch-be/internal/infrastructure/config"
)

type Server struct {
	wg     *sync.WaitGroup
	config config.HTTP
	routes map[string]chi.Router
}

func (s *Server) Serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.config.Port),
		Handler:      s.handler(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  time.Minute,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		sig := <-quit
		slog.Info("caught signal: ", "sig", sig.String())

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		slog.Info("completing background tasks")

		s.wg.Wait()
		shutdownError <- nil
	}()

	slog.Info("starting server...")

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	slog.Info("stopped server", "addr", srv.Addr)

	return nil
}

func (s *Server) AddRoutes(path string, routes chi.Router) *Server {
	s.routes[path] = routes
	return s
}

func New(wg *sync.WaitGroup, c config.HTTP) *Server {
	return &Server{
		wg:     wg,
		config: c,
		routes: make(map[string]chi.Router),
	}
}

func (s *Server) handler() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	for path, routes := range s.routes {
		r.Mount(path, routes)
	}

	return r
}
