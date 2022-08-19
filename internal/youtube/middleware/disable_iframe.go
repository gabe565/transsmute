package middleware

import (
	"context"
	"net/http"
	"strconv"
)

func DisableIframe(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()

		var disableIframe bool
		if params.Has("disable_iframe") {
			var err error
			disableIframe, err = strconv.ParseBool(params.Get("disable_iframe"))
			if err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				panic(err)
			}
		}

		r = r.WithContext(context.WithValue(r.Context(), DisableIframeKey, disableIframe))

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
