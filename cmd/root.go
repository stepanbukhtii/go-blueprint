package cmd

import (
	"context"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/stepanbukhtii/easy-tools/elog"

	"github.com/stepanbukhtii/go-blueprint/cmd/consumer"
	"github.com/stepanbukhtii/go-blueprint/cmd/server"
)

func init() {
	rootCmd.AddCommand(server.Cmd)
	rootCmd.AddCommand(consumer.Cmd)
}

var rootCmd = &cobra.Command{
	Use: "blueprint",
}

func ExecuteContext(ctx context.Context) {
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		slog.With(elog.Err(err)).Error("command execution failed")
	}
}
