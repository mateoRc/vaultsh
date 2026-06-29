package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/mateom/vaultsh/internal/httpapi"
	"github.com/mateom/vaultsh/internal/shell"
)

func main() {
	engine := shell.New()

	server := &http.Server{
		Addr:    ":8080",
		Handler: httpapi.NewHandler(engine),
	}

	log.Printf("vaultsh listening on %s", server.Addr)

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
			log.Fatal(err)
		}
		return
	case <-signalContext.Done():
		log.Print("shutting down vaultsh")
	}

	shutdownContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownContext); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}
