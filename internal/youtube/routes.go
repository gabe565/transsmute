package youtube

import (
	"context"

	"github.com/gabe565/transsmute/internal/config"
	"github.com/gabe565/transsmute/internal/youtube/channel"
	"github.com/gabe565/transsmute/internal/youtube/middleware"
	"github.com/gabe565/transsmute/internal/youtube/playlist"
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
			r.Use(middleware.IframeParam)
			r.Get("/youtube/channel/{id}", channel.Handler(service))
			r.Get("/youtube/playlist/{id}", playlist.Handler(service))
		})
	}
	return nil
}
