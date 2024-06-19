package kemono

import (
	"github.com/gabe565/transsmute/internal/config"
	"github.com/go-chi/chi/v5"
)

func Routes(r chi.Router, conf config.Kemono) {
	if !conf.Enabled {
		return
	}

	if len(conf.Hosts) == 0 {
		panic("kemono.hosts is empty.\n" +
			"If configured using an env, the value should be a JSON object with the key being the URL prefix and the value being the hostname.\n" +
			`For example: {"kemono":"kemono.su"}`)
	}

	for name, host := range conf.Hosts {
		r.Get("/"+name+"/{service}/user/{creator}", postHandler(host))
	}
}
