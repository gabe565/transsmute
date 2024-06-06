package kemono

import (
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
