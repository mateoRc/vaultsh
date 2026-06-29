package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte("ok"))
	})
	mux.Handle(
		"GET /testui/",
		http.StripPrefix("/testui/", http.FileServer(http.Dir("testui"))),
	)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Printf("vaultsh listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
