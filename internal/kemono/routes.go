package kemono

import (
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
)

func Routes(r chi.Router) {
	if !viper.GetBool("kemono.enabled") {
		return
	}

	hosts := viper.GetStringMapString("kemono.hosts")
	if len(hosts) == 0 {
		hostsRaw := viper.GetString("kemono.hosts")
		for _, v := range strings.Split(hostsRaw, ",") {
			name, val, ok := strings.Cut(v, "=")
			if ok {
				hosts[name] = val
			}
		}
	}
	if len(hosts) == 0 {
		panic("kemono.hosts is empty.\n" +
			"If configured using an env, the value should be a JSON object with the key being the URL prefix and the value being the hostname.\n" +
			`For example: {"kemono":"kemono.su"}`)
	}

	for name, host := range hosts {
		r.Get("/"+name+"/{service}/user/{creator}", postHandler(host))
	}
}
