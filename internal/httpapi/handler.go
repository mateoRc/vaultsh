package httpapi

import (
	"encoding/json"
	"net/http"
)

type execRequest struct {
	Line string `json:"line"`
}

type execResponse struct {
	Output   string `json:"output"`
	ExitCode int    `json:"exit_code"`
}

func NewHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", health)
	mux.Handle(
		"GET /testui/",
		http.StripPrefix("/testui/", http.FileServer(http.Dir("testui"))),
	)
	mux.HandleFunc("POST /api/exec", exec)

	return mux
}

func health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("ok"))
}

func exec(w http.ResponseWriter, r *http.Request) {
	var request execRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	response := execResponse{
		Output:   "Available commands:\n  help",
		ExitCode: 0,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}
