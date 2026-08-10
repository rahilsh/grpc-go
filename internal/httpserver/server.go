// Package httpserver contains the HTTP server, its routes and handlers.
package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"
)

// Config holds the runtime configuration for the HTTP server.
type Config struct {
	// Addr is the TCP address to listen on (e.g. ":3000").
	Addr string
	// Logger is used for structured logging. If nil, slog.Default() is used.
	Logger *slog.Logger
}

// Book represents a single book resource.
type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

// BooksResponse is the payload returned by the /books endpoint.
type BooksResponse struct {
	Total int    `json:"total"`
	Count int    `json:"count"`
	Books []Book `json:"books"`
}

// New builds a configured *http.Server with sane timeouts and routes.
func New(cfg Config) *http.Server {
	logger := cfg.Logger
	if logger == nil {
		logger = slog.Default()
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", rootHandler)
	mux.HandleFunc("GET /books", booksHandler)

	return &http.Server{
		Addr:              cfg.Addr,
		Handler:           logMiddleware(logger, mux),
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}
}

// Run starts srv and blocks until ctx is cancelled, then shuts it down
// gracefully. It returns any error other than http.ErrServerClosed.
func Run(ctx context.Context, srv *http.Server, logger *slog.Logger) error {
	if logger == nil {
		logger = slog.Default()
	}

	errCh := make(chan error, 1)
	go func() {
		logger.Info("server listening", slog.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		logger.Info("shutdown signal received")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	}
}

func rootHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("Hello world!"))
}

func booksHandler(w http.ResponseWriter, r *http.Request) {
	resp := BooksResponse{
		Total: 0,
		Count: 0,
		Books: []Book{},
	}
	writeJSON(r.Context(), w, http.StatusOK, resp)
}

// writeJSON encodes v as JSON to w with the given status code.
func writeJSON(_ context.Context, w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("failed to encode response", slog.String("error", err.Error()))
	}
}

// logMiddleware logs each request with its method, path, status and duration.
func logMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		logger.Info("request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", rec.status),
			slog.Duration("duration", time.Since(start)),
		)
	})
}

// statusRecorder captures the response status code for logging.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}
