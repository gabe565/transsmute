package twitter

import (
	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
)

func Routes(r chi.Router) {
	if viper.GetBool("twitter.enabled") {
		r.Get("/twitter/user/{username}", Handler)
	}
}
