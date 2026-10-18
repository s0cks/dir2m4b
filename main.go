package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"charm.land/log/v2"
	"github.com/s0cks/dir2m4b/cmd"
)

func main() {
	os.Exit(run())
}

func run() int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := cmd.RootCommand.ExecuteContext(ctx); err != nil {
		log.Error("error", "error", err)
		stop()

		if ctx.Err() != nil || (errors.Is(err, context.Canceled)) {
			return 130
		}

		return 1
	}

	return 0
}
