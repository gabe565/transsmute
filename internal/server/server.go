package server

import (
	"github.com/gabe565/transsmute/internal/docker"
	"github.com/gabe565/transsmute/internal/feed"
	"github.com/gabe565/transsmute/internal/twitter"
	"github.com/gabe565/transsmute/internal/youtube"
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
	for prefix, registerFunc := range Routers() {
		r.Route("/"+prefix, func(r chi.Router) {
			registerFunc(r)
		})
	}

	return r
}

type RoutesFunc func(r chi.Router)

func Routers() map[string]RoutesFunc {
	return map[string]RoutesFunc{
		"docker":  docker.Routes,
		"twitter": twitter.Routes,
		"youtube": youtube.Routes,
	}
}
