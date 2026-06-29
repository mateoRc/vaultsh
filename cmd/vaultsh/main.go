package main

import (
	"log"
	"net/http"

	"github.com/mateom/vaultsh/internal/httpapi"
)

func main() {
	server := &http.Server{
		Addr:    ":8080",
		Handler: httpapi.NewHandler(),
	}

	log.Printf("vaultsh listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
