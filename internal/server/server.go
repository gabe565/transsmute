package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/gabe565/transsmute/assets"
	"github.com/gabe565/transsmute/internal/config"
	"github.com/gabe565/transsmute/internal/docker"
	"github.com/gabe565/transsmute/internal/feed"
	"github.com/gabe565/transsmute/internal/kemono"
	"github.com/gabe565/transsmute/internal/youtube"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/sync/errgroup"
)

func NewServeMux(conf *config.Config) (*chi.Mux, error) {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Heartbeat("/api/health"))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Handle("/", http.NotFoundHandler())
	r.Handle("/*", http.FileServerFS(assets.Assets))

	var err error
	r.Group(func(r chi.Router) {
		r.Use(feed.SetType)
		err = errors.Join(
			docker.Routes(r, conf.Docker),
			youtube.Routes(r, conf.YouTube),
			kemono.Routes(r, conf.Kemono),
		)
	})

	return r, err
}

func ListenAndServe(ctx context.Context, conf *config.Config) error {
	mux, err := NewServeMux(conf)
	if err != nil {
		return err
	}

	group, ctx := errgroup.WithContext(ctx)

	server := http.Server{
		Addr:        conf.ListenAddress,
		Handler:     mux,
		ReadTimeout: 3 * time.Second,
	}
	group.Go(func() error {
		slog.Info("Starting server", "address", conf.ListenAddress)
		return server.ListenAndServe()
	})

	group.Go(func() error {
		<-ctx.Done()
		slog.Info("Gracefully shutting down server")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer shutdownCancel()

		return server.Shutdown(shutdownCtx)
	})

	if err := group.Wait(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
