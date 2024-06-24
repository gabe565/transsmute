package middleware

import (
	"context"
	"net/http"
)

func IframeParam(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		iframe := !params.Has("no_iframe") && !params.Has("disable_iframe")
		r = r.WithContext(context.WithValue(r.Context(), IframeKey, iframe))
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
