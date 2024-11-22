package kemono

import (
	"errors"
	"net/http"
	"net/url"
	"path"

	"github.com/go-chi/chi/v5"
)

func redirectHandler(host string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := &url.URL{RawQuery: r.URL.RawQuery}
		if id := chi.URLParam(r, "id"); id != "" {
			u.Path = path.Join("id", id)
		} else {
			service := chi.URLParam(r, "service")
			username := chi.URLParam(r, "name")
			creator, err := GetCreatorByUsername(r.Context(), host, service, username)
			if err != nil {
				var respErr UpstreamResponseError
				if errors.As(err, &respErr) {
					http.Error(w, respErr.Body(), respErr.Response.StatusCode)
					return
				} else if errors.Is(err, ErrCreatorNotFound) {
					http.Error(w, err.Error(), http.StatusNotFound)
					return
				}
				panic(err)
			}
			u.Path = path.Join("../id", creator.ID)
		}
		http.Redirect(w, r, u.String(), http.StatusPermanentRedirect)
	}
}
