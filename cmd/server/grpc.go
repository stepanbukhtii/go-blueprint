package server

import (
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/stepanbukhtii/easy-tools/elog"

	"github.com/stepanbukhtii/go-blueprint/internal/app"
	internalgrpc "github.com/stepanbukhtii/go-blueprint/internal/transport/grpc"
)

var grpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Start GRPC server",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		app, err := app.New(ctx)
		if err != nil {
			return err
		}

		grpcServer := internalgrpc.NewServer(app)

		go func() {
			slog.InfoContext(ctx, "grpc server started")

			if err := grpcServer.Serve(app.Config.GRPC); err != nil {
				slog.With(elog.Err(err)).ErrorContext(ctx, "grpc server running failed")
			}
		}()

		<-ctx.Done()

		slog.Info("received termination signal, shutting down gracefully...")

		grpcServer.GracefulStop()

		app.Close()

		slog.Info("grpc server shutdown finished")

		return nil
	},
}
