package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/stepanbukhtii/go-blueprint/cmd"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cmd.ExecuteContext(ctx)
}
