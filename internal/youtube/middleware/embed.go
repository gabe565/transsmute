package middleware

import (
	"net/http"
)

func EmbedParam(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		embed := !params.Has("no_embed") && !params.Has("no_iframe") && !params.Has("disable_iframe")
		r = r.WithContext(NewEmbedContext(r.Context(), embed))
		next.ServeHTTP(w, r)
	})
}
