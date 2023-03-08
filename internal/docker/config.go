package docker

import (
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	flag.Bool("docker-enabled", true, "Docker API enabled")
	if err := viper.BindPFlag("docker.enabled", flag.Lookup("docker-enabled")); err != nil {
		panic(err)
	}
}
