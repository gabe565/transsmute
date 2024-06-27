package kemono

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gabe565/transsmute/internal/feed"
	"github.com/gabe565/transsmute/internal/util"
	"github.com/go-chi/chi/v5"
)

func podcastHandler(host string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		creator, err := GetCreatorInfo(r.Context(), host, chi.URLParam(r, "service"), chi.URLParam(r, "creator"))
		if err != nil {
			if errors.Is(err, ErrCreatorNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			} else if errors.Is(err, util.ErrUpstreamRequest) {
				http.Error(w, err.Error(), http.StatusServiceUnavailable)
				return
			}
			panic(err)
		}

		pagesRaw := r.URL.Query().Get("pages")
		pages := uint64(1)
		if pagesRaw != "" {
			var err error
			if pages, err = strconv.ParseUint(pagesRaw, 10, 64); err != nil || pages == 0 {
				http.Error(w, "pages must be a positive integer", http.StatusBadRequest)
				return
			}
		}

		tag := r.URL.Query().Get("tag")
		query := r.URL.Query().Get("q")

		f, err := creator.Podcast(r.Context(), pages, tag, query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			panic(err)
		}

		if val := r.URL.Query().Get("title"); val != "" {
			f.Title = val
			if f.Image != nil {
				f.Image.Title = val
			}
		}

		r = r.WithContext(context.WithValue(r.Context(), feed.FeedKey, f))
		if err := feed.WriteFeed(w, r); err != nil {
			panic(err)
		}
	}
}
