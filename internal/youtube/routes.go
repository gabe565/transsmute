package youtube

import (
	"github.com/gabe565/tuberss/internal/youtube/channel"
	"github.com/gabe565/tuberss/internal/youtube/middleware"
	"github.com/gabe565/tuberss/internal/youtube/playlist"
	"github.com/go-chi/chi/v5"
)

func Routes(r chi.Router, prefix string) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.DisableIframe)

		r.Get("/"+prefix+"/channel/{id}", channel.Handler)
		r.Get("/"+prefix+"/playlist/{id}", playlist.Handler)
	})
}
