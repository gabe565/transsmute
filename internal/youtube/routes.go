package youtube

import (
	"context"

	"gabe565.com/transsmute/internal/config"
	"gabe565.com/transsmute/internal/youtube/channel"
	"gabe565.com/transsmute/internal/youtube/middleware"
	"gabe565.com/transsmute/internal/youtube/playlist"
	"github.com/go-chi/chi/v5"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func Routes(r chi.Router, conf config.YouTube) error {
	if conf.Enabled {
		service, err := youtube.NewService(
			context.Background(),
			option.WithAPIKey(conf.APIKey),
			option.WithTelemetryDisabled(),
		)
		if err != nil {
			return err
		}

		r.Group(func(r chi.Router) {
			r.Use(middleware.EmbedParam, middleware.LimitParam)
			r.Get("/youtube/channel/id/{id}", channel.Handler(service))
			r.Get("/youtube/channel/username/{username}", channel.RedirectHandler(service))
			r.Get("/youtube/channel/handle/{handle}", channel.RedirectHandler(service))
			r.Get("/youtube/playlist/{id}", playlist.Handler(service))

			r.Get("/youtube/channel/{id}", channel.RedirectHandler(service)) // Deprecated
		})
	}
	return nil
}
