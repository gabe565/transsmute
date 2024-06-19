package docker

import (
	"github.com/gabe565/transsmute/internal/config"
	"github.com/go-chi/chi/v5"
)

func Routes(r chi.Router, conf config.Docker) error {
	if conf.Enabled {
		if err := SetupRegistries(conf); err != nil {
			return err
		}
		r.Get("/docker/tags/*", Handler)
	}
	return nil
}
