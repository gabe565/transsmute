package main

import (
	"github.com/gabe565/tuberss/internal/server"
	flag "github.com/spf13/pflag"
	"log"
	"net/http"
)

func main() {
	var address string
	flag.StringVar(&address, "address", ":3000", "Listening address")

	flag.Parse()

	s := server.New()
	log.Println("Listening on " + address)
	if err := http.ListenAndServe(address, s.Handler()); err != nil {
		panic(err)
	}
}
