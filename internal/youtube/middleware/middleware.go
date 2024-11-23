package middleware

import (
	"net/http"
	"strconv"
)

func EmbedParam(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		embed := !params.Has("no_embed") && !params.Has("no_iframe") && !params.Has("disable_iframe")
		r = r.WithContext(NewEmbedContext(r.Context(), embed))
		next.ServeHTTP(w, r)
	})
}

const DefaultLimit = 15

func LimitParam(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limit := DefaultLimit
		if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
			val, err := strconv.Atoi(limitStr)
			if err != nil {
				http.Error(w, "limit must be an integer", http.StatusBadRequest)
				return
			}
			if val != 0 {
				limit = val
			}
		}
		r = r.WithContext(NewLimitContext(r.Context(), limit))
		next.ServeHTTP(w, r)
	})
}
