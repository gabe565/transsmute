package playlist

import (
	"context"
	"errors"
	"net/http"

	"github.com/gabe565/transsmute/internal/feed"
	"github.com/gabe565/transsmute/internal/youtube/middleware"
	"github.com/go-chi/chi/v5"
	"google.golang.org/api/youtube/v3"
)

func Handler(service *youtube.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		identifier := chi.URLParam(r, "id")
		plist := New(service, identifier)

		f, err := plist.Feed(r.Context(), r.Context().Value(middleware.NoIframeKey).(bool))
		if err != nil {
			if errors.Is(err, ErrInvalid) {
				http.Error(w, "404 playlist not found", http.StatusNotFound)
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
