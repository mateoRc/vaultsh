package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/mateom/vaultsh/internal/external"
	"github.com/mateom/vaultsh/internal/httpapi"
	"github.com/mateom/vaultsh/internal/shell"
	"github.com/mateom/vaultsh/internal/storage"
	"github.com/mateom/vaultsh/internal/telemetry"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	contentPath := os.Getenv("CONTENT_PATH")
	if contentPath == "" {
		contentPath = "/app/content"
	}

	root, err := storage.Load(os.DirFS(contentPath))
	if err != nil {
		logger.Error("content loading failed", "error", err)
		os.Exit(1)
	}
	atlasURL := os.Getenv("ATLAS_URL")
	forgeURL := os.Getenv("FORGE_URL")
	services := external.NewClient(
		atlasURL,
		forgeURL,
		serviceToken(atlasURL, "ATLAS_AUTH_TOKEN", logger),
		serviceToken(forgeURL, "FORGE_AUTH_TOKEN", logger),
	)
	events := telemetry.NewDispatcher(services, telemetryQueueSize(), logger)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		events.Close(ctx)
	}()
	sessions := shell.NewSessionManagerWithConfigAndDependencies(
		root,
		shell.SessionConfig{MaxSessions: sessionLimit()},
		shell.Dependencies{
			Search:  services,
			Metrics: services,
			Events:  events,
		},
	)

	server := &http.Server{
		Addr: ":8080",
		Handler: httpapi.NewHandlerWithConfig(
			sessions,
			services,
			httpapi.HandlerConfig{TrustProxyHeaders: trustProxyHeaders()},
		),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    16 * 1024,
	}

	logger.Info("server started", "address", server.Addr)

	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- server.ListenAndServe()
	}()

	signalContext, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	go sessions.RunCleanup(signalContext, shell.DefaultSessionCleanupInterval)

	select {
	case err := <-serverErrors:
		if err != nil && err != http.ErrServerClosed {
			logger.Error("server failed", "error", err)
			os.Exit(1)
		}
		return
	case <-signalContext.Done():
		logger.Info("server stopping")
	}

	shutdownContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownContext); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
	}
}

func serviceToken(serviceURL, name string, logger *slog.Logger) string {
	if serviceURL == "" {
		return ""
	}
	value := os.Getenv(name)
	if value == "" {
		logger.Error("required environment variable is missing", "name", name)
		os.Exit(1)
	}
	return value
}

func telemetryQueueSize() int {
	value, err := strconv.Atoi(os.Getenv("TELEMETRY_QUEUE_SIZE"))
	if err != nil || value <= 0 {
		return 1000
	}
	return value
}

func sessionLimit() int {
	value, err := strconv.Atoi(os.Getenv("SESSION_LIMIT"))
	if err != nil || value <= 0 {
		return shell.DefaultMaxSessions
	}
	return value
}

func trustProxyHeaders() bool {
	value, err := strconv.ParseBool(os.Getenv("TRUST_PROXY_HEADERS"))
	return err == nil && value
}
