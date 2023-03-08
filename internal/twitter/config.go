package twitter

import (
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	flag.Bool("twitter-enabled", true, "Twitter API enabled")
	if err := viper.BindPFlag("twitter.enabled", flag.Lookup("twitter-enabled")); err != nil {
		panic(err)
	}
}
