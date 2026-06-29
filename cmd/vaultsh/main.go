package main

import (
	"encoding/json"
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
	mux.HandleFunc("POST /api/exec", func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Line string `json:"line"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(struct {
			Output   string `json:"output"`
			ExitCode int    `json:"exit_code"`
		}{
			Output:   "Available commands:\n  help",
			ExitCode: 0,
		})
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Printf("vaultsh listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
