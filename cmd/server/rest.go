package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/spf13/cobra"
	"github.com/stepanbukhtii/easy-tools/elog"

	"github.com/stepanbukhtii/go-blueprint/internal/app"
	internalhttp "github.com/stepanbukhtii/go-blueprint/internal/transport/http"
)

const shutdownTimeout = 30 * time.Second

var restCmd = &cobra.Command{
	Use:   "rest",
	Short: "Start REST server",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		app, err := app.New(ctx)
		if err != nil {
			return err
		}

		restServer, err := internalhttp.NewServer(app)
		if err != nil {
			return err
		}

		go func() {
			slog.InfoContext(ctx, "rest server started")
			if err := restServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				slog.With(elog.Err(err)).ErrorContext(ctx, "rest server running failed")
			}
		}()

		<-ctx.Done()

		slog.Info("received termination signal, shutting down gracefully...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := restServer.Shutdown(shutdownCtx); err != nil {
			slog.With(elog.Err(err)).Error("rest server shutdown failed")
		}

		app.Close()

		slog.Info("rest server shutdown finished")

		return nil
	},
}
