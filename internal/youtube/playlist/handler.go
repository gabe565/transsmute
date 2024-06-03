package playlist

import (
	"context"
	"errors"
	"net/http"

	"github.com/gabe565/transsmute/internal/feed"
	"github.com/gabe565/transsmute/internal/youtube/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	service, err := youtube.NewService(
		r.Context(),
		option.WithAPIKey(viper.GetString("youtube.key")),
		option.WithTelemetryDisabled(),
	)
	if err != nil {
		panic(err)
	}

	identifier := chi.URLParam(r, "id")
	plist := New(service, r.Context(), identifier)

	f, err := plist.Feed(r.Context().Value(middleware.DisableIframeKey).(bool))
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
