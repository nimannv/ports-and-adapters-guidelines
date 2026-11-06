// Composition root: config -> outbound adapters -> use cases -> inbound adapters.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	// 1. load config  2. build outbound adapters  3. build use cases  4. build and start inbound adapters
	<-ctx.Done()
	return nil
}
