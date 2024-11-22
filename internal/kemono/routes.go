package kemono

import (
	"errors"

	"gabe565.com/transsmute/internal/config"
	"github.com/go-chi/chi/v5"
)

var ErrNoHosts = errors.New("kemono.hosts is empty")

func Routes(r chi.Router, conf config.Kemono) error {
	if conf.Enabled {
		if len(conf.Hosts) == 0 {
			return ErrNoHosts
		}

		for name, host := range conf.Hosts {
			r.Get("/"+name+"/{service}/user/name/{name}", redirectHandler(host))
			r.Get("/"+name+"/{service}/user/id/{id}", postHandler(host))
			r.Get("/"+name+"/{service}/podcast/name/{name}", redirectHandler(host))
			r.Get("/"+name+"/{service}/podcast/id/{id}", podcastHandler(host))

			r.Get("/"+name+"/{service}/user/{id}", redirectHandler(host))    // Deprecated
			r.Get("/"+name+"/{service}/podcast/{id}", redirectHandler(host)) // Deprecated
		}
	}
	return nil
}
