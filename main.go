package main

import (
	"os"

	"gabe565.com/transsmute/cmd"
	"gabe565.com/utils/cobrax"
)

var version = "beta"

func main() {
	root := cmd.New(cobrax.WithVersion(version))
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
