package kemono

import (
	"errors"
	"net/http"
	"strconv"

	"gabe565.com/transsmute/internal/feed"
	"github.com/go-chi/chi/v5"
)

func podcastHandler(host string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		creator, err := GetCreatorByID(r.Context(), host, chi.URLParam(r, "service"), chi.URLParam(r, "id"))
		if err != nil {
			var respErr UpstreamResponseError
			if errors.As(err, &respErr) {
				http.Error(w, respErr.Body(), respErr.Response.StatusCode)
				return
			}
			panic(err)
		}

		pages := uint64(1)
		if v := r.URL.Query().Get("pages"); v != "" {
			var err error
			if pages, err = strconv.ParseUint(v, 10, 64); err != nil || pages == 0 {
				http.Error(w, "pages must be a positive integer", http.StatusBadRequest)
				return
			}
		}

		tag := r.URL.Query().Get("tag")
		query := r.URL.Query().Get("q")

		f, err := creator.Podcast(r.Context(), pages, tag, query)
		if err != nil {
			var respErr UpstreamResponseError
			if errors.As(err, &respErr) {
				http.Error(w, respErr.Body(), respErr.Response.StatusCode)
				return
			}
			panic(err)
		}

		if val := r.URL.Query().Get("title"); val != "" {
			f.Title = val
			if f.Image != nil {
				f.Image.Title = val
			}
		}

		if err := feed.WritePodcast(w, r, f); err != nil {
			panic(err)
		}
	}
}
