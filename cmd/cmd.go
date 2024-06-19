package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gabe565/transsmute/internal/docker"
	"github.com/gabe565/transsmute/internal/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Short: "transsmute",
		RunE:  run,
	}

	cmd.Flags().String("address", ":3000", "Listening address")
	if err := viper.BindPFlag("address", cmd.Flags().Lookup("address")); err != nil {
		panic(err)
	}

	return cmd
}

func run(_ *cobra.Command, _ []string) error {
	initViper()

	if err := docker.SetupRegistries(); err != nil {
		return err
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	if err := server.ListenAndServe(ctx, viper.GetString("address")); err != nil {
		return err
	}
	return nil
}
