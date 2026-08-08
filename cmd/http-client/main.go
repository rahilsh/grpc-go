// Command http-client performs a single HTTP GET and prints the response body.
package main

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	url := os.Getenv("CLIENT_URL")
	if url == "" {
		url = "https://jsonplaceholder.typicode.com/todos/1"
	}

	if err := run(url, logger); err != nil {
		logger.Error("request failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func run(url string, logger *slog.Logger) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logger.Error("failed to close response body", slog.String("error", err.Error()))
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	logger.Info("response received",
		slog.Int("status", resp.StatusCode),
		slog.String("body", string(body)),
	)
	return nil
}
