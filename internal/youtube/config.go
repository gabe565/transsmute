package youtube

import (
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	flag.Bool("youtube-enabled", true, "YouTube API enabled")
	if err := viper.BindPFlag("youtube.enabled", flag.Lookup("youtube-enabled")); err != nil {
		panic(err)
	}

	flag.String("youtube-key", "", "YouTube API key")
	if err := viper.BindPFlag("youtube.key", flag.Lookup("youtube-key")); err != nil {
		panic(err)
	}
}
