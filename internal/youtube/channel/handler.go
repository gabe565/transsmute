package channel

import (
	"context"
	"errors"
	"github.com/gabe565/tuberss/internal/feed"
	"github.com/gabe565/tuberss/internal/youtube/config"
	"github.com/gabe565/tuberss/internal/youtube/middleware"
	"github.com/gabe565/tuberss/internal/youtube/playlist"
	"github.com/go-chi/chi/v5"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	service, err := youtube.NewService(r.Context(), option.WithAPIKey(config.ApiKey))
	if err != nil {
		panic(err)
	}

	identifier := chi.URLParam(r, "id")
	ch := New(service, r.Context(), identifier)

	f, err := ch.Feed(r.Context().Value(middleware.DisableIframeKey).(bool))
	if err != nil {
		if errors.Is(err, ErrInvalid) {
			http.Error(w, "404 channel not found", http.StatusNotFound)
			return
		} else if errors.Is(err, playlist.ErrInvalid) {
			http.Error(w, "404 channel has no videos", http.StatusNotFound)
			return
		}
		panic(err)
	}

	r = r.WithContext(context.WithValue(r.Context(), feed.FeedKey, f))

	if err := feed.WriteFeed(w, r); err != nil {
		panic(err)
	}
}
