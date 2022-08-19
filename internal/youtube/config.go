package youtube

import (
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	flag.String("youtube-key", "", "YouTube API key")
	if err := viper.BindPFlag("youtube.key", flag.Lookup("youtube-key")); err != nil {
		panic(err)
	}
}
