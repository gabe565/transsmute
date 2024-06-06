package main

import (
	"log"
	"net/http"
	"strings"
	"time"

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
		//nolint:errorlint
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Config file was found but another error was produced
			panic("fatal error reading config file:" + err.Error())
		}
	}

	flag.Parse()

	if err := docker.SetupRegistries(); err != nil {
		panic(err)
	}

	s := server.New()
	address := viper.GetString("address")
	log.Println("Listening on " + address)

	srv := http.Server{
		Addr:        address,
		Handler:     s.Handler(),
		ReadTimeout: 3 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
