package cmd

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"gabe565.com/transsmute/internal/config"
	"gabe565.com/transsmute/internal/server"
	"gabe565.com/utils/cobrax"
	"github.com/spf13/cobra"
)

func New(opts ...cobrax.Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "transsmute",
		Long: "Build RSS feeds for websites that don't provide them.",
		RunE: run,

		DisableAutoGenTag: true,
	}
	conf := config.New()
	conf.RegisterFlags(cmd)
	cmd.SetContext(config.NewContext(context.Background(), conf))

	for _, opt := range opts {
		opt(cmd)
	}

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

	slog.Info("Transsmute", "version", cobrax.GetVersion(cmd), "commit", cobrax.GetCommit(cmd))

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	if err := server.ListenAndServe(ctx, conf); err != nil {
		return err
	}
	return nil
}
