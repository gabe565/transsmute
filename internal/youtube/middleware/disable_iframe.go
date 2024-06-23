package middleware

import (
	"context"
	"net/http"
)

func NoIframe(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		noIframe := params.Has("no_iframe") || params.Has("disable_iframe")
		r = r.WithContext(context.WithValue(r.Context(), NoIframeKey, noIframe))
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
