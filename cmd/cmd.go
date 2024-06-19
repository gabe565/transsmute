package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gabe565/transsmute/internal/config"
	"github.com/gabe565/transsmute/internal/server"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Short: "transsmute",
		RunE:  run,
	}
	conf := config.New()
	conf.RegisterFlags(cmd)
	cmd.SetContext(config.NewContext(context.Background(), conf))
	return cmd
}

func run(cmd *cobra.Command, _ []string) error {
	conf, ok := config.FromContext(cmd.Context())
	if !ok {
		panic("config missing from command context")
	}

	if err := conf.Load(cmd); err != nil {
		return err
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	if err := server.ListenAndServe(ctx, conf); err != nil {
		return err
	}
	return nil
}
