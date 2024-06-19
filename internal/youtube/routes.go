package youtube

import (
	"github.com/gabe565/transsmute/internal/config"
	"github.com/gabe565/transsmute/internal/youtube/channel"
	"github.com/gabe565/transsmute/internal/youtube/middleware"
	"github.com/gabe565/transsmute/internal/youtube/playlist"
	"github.com/go-chi/chi/v5"
)

func Routes(r chi.Router, conf config.YouTube) {
	if conf.Enabled {
		r.Group(func(r chi.Router) {
			r.Use(middleware.DisableIframe)

			r.Get("/youtube/channel/{id}", channel.Handler(conf.APIKey))
			r.Get("/youtube/playlist/{id}", playlist.Handler(conf.APIKey))
		})
	}
}
