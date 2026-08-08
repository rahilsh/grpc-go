// Command http-server starts the HTTP API.
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/rahilsh/golang-lab/internal/httpserver"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	addr := os.Getenv("SERVER_ADDR")
	if addr == "" {
		addr = ":3000"
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv := httpserver.New(httpserver.Config{Addr: addr, Logger: logger})

	if err := httpserver.Run(ctx, srv, logger); err != nil {
		logger.Error("server error", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.Info("server stopped cleanly")
}
