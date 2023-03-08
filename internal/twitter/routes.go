package twitter

import (
	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
)

func Routes(r chi.Router, prefix string) {
	if viper.GetBool("twitter.enabled") {
		r.Get("/"+prefix+"/user/{username}", Handler)
	}
}
