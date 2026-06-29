package main

import (
	"log"
	"net/http"

	"github.com/mateom/vaultsh/internal/httpapi"
	"github.com/mateom/vaultsh/internal/shell"
)

func main() {
	engine := &shell.Engine{}

	server := &http.Server{
		Addr:    ":8080",
		Handler: httpapi.NewHandler(engine),
	}

	log.Printf("vaultsh listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
