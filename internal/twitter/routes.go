package twitter

import "github.com/go-chi/chi/v5"

func Routes(r chi.Router, prefix string) {
	r.Get("/"+prefix+"/user/{username}", Handler)
}
