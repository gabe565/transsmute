package docker

import (
	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
)

func Routes(r chi.Router) {
	if viper.GetBool("docker.enabled") {
		r.Get("/docker/tags/*", Handler)
	}
}
