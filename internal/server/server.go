package server

import (
	"net/http"

	"github.com/gabe565/transsmute/assets"
	"github.com/gabe565/transsmute/internal/docker"
	"github.com/gabe565/transsmute/internal/feed"
	"github.com/gabe565/transsmute/internal/kemono"
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

	r.Get("/*", http.FileServer(http.FS(assets.Assets)).ServeHTTP)

	docker.Routes(r)
	youtube.Routes(r)
	kemono.Routes(r)

	return r
}
