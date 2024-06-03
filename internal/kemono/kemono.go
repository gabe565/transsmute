package kemono

import (
	"net/url"
	"path"

	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	flag.Bool("kemono-enabled", true, "Kemono API enabled")
	if err := viper.BindPFlag("kemono.enabled", flag.Lookup("kemono-enabled")); err != nil {
		panic(err)
	}

	flag.StringToString("kemono-hosts", map[string]string{"kemono": "kemono.su"}, "Kemono API hosts, where the key is the URL prefix and the value is the host")
	if err := viper.BindPFlag("kemono.hosts", flag.Lookup("kemono-hosts")); err != nil {
		panic(err)
	}
}

func publicURL(host string, c Creator) *url.URL {
	return &url.URL{
		Scheme: "https",
		Host:   host,
		Path:   path.Join(c.Service, "user", c.ID),
	}
}

func postURL(host string, c Creator, p Post) *url.URL {
	u := publicURL(host, c)
	u.Path = path.Join(u.Path, "post", p.ID)
	return u
}

func postAPIURL(host string, c Creator) *url.URL {
	return &url.URL{
		Scheme: "https",
		Host:   host,
		Path:   path.Join("api/v1", c.Service, "user", c.ID),
	}
}
