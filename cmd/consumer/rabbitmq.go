package consumer

import (
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/stepanbukhtii/easy-tools/elog"

	"github.com/stepanbukhtii/go-blueprint/internal/app"
	"github.com/stepanbukhtii/go-blueprint/internal/transport/rabbitmq"
)

var rabbitMQCmd = &cobra.Command{
	Use:   "rabbitmq",
	Short: "Start rabbitMQ consumer",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		app, err := app.New(ctx)
		if err != nil {
			return err
		}

		rabbitMQConsumerGroup, err := rabbitmq.NewConsumers(app)
		if err != nil {
			return err
		}

		go func() {
			slog.InfoContext(ctx, "rabbitmq consumer started")

			if err := rabbitMQConsumerGroup.Consume(ctx); err != nil {
				slog.With(elog.Err(err)).ErrorContext(ctx, "rabbitmq consuming failed")
			}
		}()

		<-ctx.Done()

		slog.Info("received termination signal, shutting down gracefully...")

		if err := rabbitMQConsumerGroup.Shutdown(); err != nil {
			slog.With(elog.Err(err)).Error("rabbitmq consumer shutdown failed")
		}

		app.Close()

		slog.Info("rabbitmq consumer shutdown finished")

		return nil
	},
}
