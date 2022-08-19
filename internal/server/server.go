package server

import (
	"github.com/gabe565/tuberss/internal/feed"
	"github.com/gabe565/tuberss/internal/youtube"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New() Server {
	return Server{}
}

type Server struct{}

func (s Server) Handler() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Heartbeat("/api/health"))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(feed.SetType)
	for name, router := range Routers() {
		router(r, name)
	}

	return r
}

type RoutesFunc func(r chi.Router, prefix string)

func Routers() map[string]RoutesFunc {
	return map[string]RoutesFunc{
		"youtube": youtube.Routes,
	}
}
