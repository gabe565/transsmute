package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gabe565/transsmute/internal/docker"
	"github.com/gabe565/transsmute/internal/server"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	flag.String("address", ":3000", "Listening address")
	if err := viper.BindPFlag("address", flag.Lookup("address")); err != nil {
		panic(err)
	}
}

func main() {
	viper.SetConfigName("transsmute")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/")
	viper.AddConfigPath("$HOME/")
	viper.AddConfigPath("/etc/transsmute/")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("transsmute")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error
		} else {
			// Config file was found but another error was produced
			panic(fmt.Errorf("Fatal error reading config file: %w \n", err))
		}
	}

	flag.Parse()

	if err := docker.SetupRegistries(); err != nil {
		panic(err)
	}

	s := server.New()
	address := viper.GetString("address")
	log.Println("Listening on " + address)
	if err := http.ListenAndServe(address, s.Handler()); err != nil {
		panic(err)
	}
}
