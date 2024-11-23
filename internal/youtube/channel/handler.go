package channel

import (
	"errors"
	"net/http"
	"net/url"
	"path"

	"gabe565.com/transsmute/internal/feed"
	"gabe565.com/transsmute/internal/youtube/middleware"
	"gabe565.com/transsmute/internal/youtube/playlist"
	"github.com/go-chi/chi/v5"
	"google.golang.org/api/youtube/v3"
)

func Handler(service *youtube.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		channelID := chi.URLParam(r, "id")
		username := chi.URLParam(r, "username")

		ch := New(service, channelID, username)
		ch.Embed = middleware.EmbedFromContext(r.Context())

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

		if err := feed.WriteFeed(w, r, f); err != nil {
			panic(err)
		}
	}
}

func RedirectHandler(service *youtube.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := &url.URL{
			Path:     "/youtube/channel/id",
			RawQuery: r.URL.RawQuery,
		}

		if id := chi.URLParam(r, "id"); id != "" {
			u.Path = path.Join(u.Path, id)
		} else {
			username := chi.URLParam(r, "username")
			handle := chi.URLParam(r, "handle")

			call := service.Channels.List([]string{"id"})
			call.Context(r.Context())

			var reqType string
			switch {
			case username != "":
				call.ForUsername(username)
				reqType = "Username"
			case handle != "":
				call.ForHandle(handle)
				reqType = "Handle"
			default:
				http.Error(w, reqType+" is required", http.StatusBadRequest)
				return
			}

			resp, err := call.Do()
			if err != nil {
				panic(err)
			}

			if len(resp.Items) == 0 {
				http.Error(w, reqType+" not found", http.StatusNotFound)
				return
			}

			u.Path = path.Join(u.Path, resp.Items[0].Id)
		}

		http.Redirect(w, r, u.String(), http.StatusPermanentRedirect)
	}
}
