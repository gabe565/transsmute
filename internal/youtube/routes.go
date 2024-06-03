package youtube

import (
	"github.com/gabe565/transsmute/internal/youtube/channel"
	"github.com/gabe565/transsmute/internal/youtube/middleware"
	"github.com/gabe565/transsmute/internal/youtube/playlist"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
)

func Routes(r chi.Router) {
	if viper.GetBool("youtube.enabled") {
		r.Group(func(r chi.Router) {
			r.Use(middleware.DisableIframe)

			r.Get("/youtube/channel/{id}", channel.Handler)
			r.Get("/youtube/playlist/{id}", playlist.Handler)
		})
	}
}
