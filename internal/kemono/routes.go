package kemono

import (
	"errors"

	"github.com/gabe565/transsmute/internal/config"
	"github.com/go-chi/chi/v5"
)

var ErrNoHosts = errors.New("kemono.hosts is empty")

func Routes(r chi.Router, conf config.Kemono) error {
	if conf.Enabled {
		if len(conf.Hosts) == 0 {
			return ErrNoHosts
		}

		initCreatorCache()

		for name, host := range conf.Hosts {
			r.Get("/"+name+"/{service}/user/{creator}", postHandler(host))
			r.Get("/"+name+"/{service}/podcast/{creator}", podcastHandler(host))
		}
	}
	return nil
}
