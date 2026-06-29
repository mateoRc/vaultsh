package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mateom/vaultsh/internal/httpapi"
	"github.com/mateom/vaultsh/internal/shell"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	engine := shell.New()

	server := &http.Server{
		Addr:    ":8080",
		Handler: httpapi.NewHandler(engine),
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
