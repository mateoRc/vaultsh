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
	services := external.NewClient(os.Getenv("ATLAS_URL"), os.Getenv("FORGE_URL"))
	events := telemetry.NewDispatcher(services, telemetryQueueSize(), logger)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		events.Close(ctx)
	}()
	sessions := shell.NewSessionManagerWithDependencies(root, shell.Dependencies{
		Search:  services,
		Metrics: services,
		Events:  events,
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: httpapi.NewHandlerWithStatus(sessions, services),
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

func telemetryQueueSize() int {
	value, err := strconv.Atoi(os.Getenv("TELEMETRY_QUEUE_SIZE"))
	if err != nil || value <= 0 {
		return 1000
	}
	return value
}
