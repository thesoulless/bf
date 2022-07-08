package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/thesoulless/bf/internal/cmd"
)

func main() {
	os.Exit(run())
}

func run() int {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	return cmd.Execute(ctx)
}
