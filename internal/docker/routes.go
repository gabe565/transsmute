package docker

import (
	"github.com/gabe565/transsmute/internal/config"
	"github.com/go-chi/chi/v5"
)

func Routes(r chi.Router, conf config.Docker) {
	if conf.Enabled {
		r.Get("/docker/tags/*", Handler)
	}
}
