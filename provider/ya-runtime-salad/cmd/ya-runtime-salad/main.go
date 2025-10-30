package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/afero"

	"github.com/seniorquico/golem-network-demo/provider/ya-runtime-salad/internal/runtime"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	os.Exit(runtime.HandleInvocation(ctx, os.Args, afero.NewOsFs(), os.Stdout, os.Stderr))
}
