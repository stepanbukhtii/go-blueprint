package consumer

import (
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/stepanbukhtii/easy-tools/elog"

	"github.com/stepanbukhtii/go-blueprint/internal/app"
	"github.com/stepanbukhtii/go-blueprint/internal/transport/nats"
)

var natsCmd = &cobra.Command{
	Use:   "nats",
	Short: "Start nats subscriber",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		app, err := app.New(ctx)
		if err != nil {
			return err
		}

		natsSubscriber, err := nats.NewSubscriber(app)
		if err != nil {
			return err
		}

		slog.InfoContext(ctx, "nats subscriber started")

		<-ctx.Done()

		slog.Info("received termination signal, shutting down gracefully...")

		if err := natsSubscriber.Shutdown(); err != nil {
			slog.With(elog.Err(err)).ErrorContext(ctx, "nats subscriber shutdown failed")
		}

		app.Close()

		slog.Info("nats subscriber shutdown finished")

		return nil
	},
}
