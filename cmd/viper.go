package cmd

import (
	"strings"

	"github.com/spf13/viper"
)

func initViper() {
	viper.SetConfigName("transsmute")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/")
	viper.AddConfigPath("$HOME/")
	viper.AddConfigPath("/etc/transsmute/")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("transsmute")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	if err := viper.ReadInConfig(); err != nil {
		//nolint:errorlint
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Config file was found but another error was produced
			panic("fatal error reading config file:" + err.Error())
		}
	}
}
