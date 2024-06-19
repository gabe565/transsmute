package main

import (
	"os"

	"github.com/gabe565/transsmute/cmd"
)

func main() {
	if err := cmd.New().Execute(); err != nil {
		os.Exit(1)
	}
}
