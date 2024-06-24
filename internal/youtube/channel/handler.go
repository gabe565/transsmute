package channel

import (
	"context"
	"errors"
	"net/http"

	"github.com/gabe565/transsmute/internal/feed"
	"github.com/gabe565/transsmute/internal/youtube/middleware"
	"github.com/gabe565/transsmute/internal/youtube/playlist"
	"github.com/go-chi/chi/v5"
	"google.golang.org/api/youtube/v3"
)

func Handler(service *youtube.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		identifier := chi.URLParam(r, "id")
		iframe := r.Context().Value(middleware.IframeKey).(bool)
		ch := New(service, identifier)
		ch.Iframe = iframe

		f, err := ch.Feed(r.Context())
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
}
