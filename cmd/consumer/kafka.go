package consumer

import (
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/stepanbukhtii/easy-tools/elog"

	"github.com/stepanbukhtii/go-blueprint/internal/app"
	"github.com/stepanbukhtii/go-blueprint/internal/transport/kafka"
)

var kafkaCmd = &cobra.Command{
	Use:   "kafka",
	Short: "Start kafka consumer",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		app, err := app.New(ctx)
		if err != nil {
			return err
		}

		kafkaConsumerGroup := kafka.NewConsumers(app)

		go func() {
			slog.InfoContext(ctx, "kafka consumer started")

			if err := kafkaConsumerGroup.Consume(ctx); err != nil {
				slog.With(elog.Err(err)).ErrorContext(ctx, "consume kafka")
			}
		}()

		<-ctx.Done()

		slog.Info("received termination signal, shutting down gracefully...")

		kafkaConsumerGroup.Shutdown()

		app.Close()

		slog.Info("kafka consumer shutdown finished")

		return nil
	},
}
