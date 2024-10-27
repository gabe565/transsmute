package playlist

import (
	"errors"
	"net/http"

	"gabe565.com/transsmute/internal/feed"
	"gabe565.com/transsmute/internal/youtube/middleware"
	"github.com/go-chi/chi/v5"
	"google.golang.org/api/youtube/v3"
)

func Handler(service *youtube.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		identifier := chi.URLParam(r, "id")
		iframe := r.Context().Value(middleware.IframeKey).(bool)
		plist := New(service, identifier)
		plist.Iframe = iframe

		f, err := plist.Feed(r.Context())
		if err != nil {
			if errors.Is(err, ErrInvalid) {
				http.Error(w, "404 playlist not found", http.StatusNotFound)
				return
			}
			panic(err)
		}

		r = r.WithContext(feed.NewContext(r.Context(), f))
		if err := feed.WriteFeed(w, r); err != nil {
			panic(err)
		}
	}
}
