package config

import flag "github.com/spf13/pflag"

var ApiKey string

func init() {
	flag.StringVar(&ApiKey, "key", "", "YouTube API key")
}
