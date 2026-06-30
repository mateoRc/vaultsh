package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/mateom/vaultsh/internal/command"
	"github.com/mateom/vaultsh/internal/shell"
)

type execRequest struct {
	Line      string `json:"line"`
	SessionID string `json:"session_id,omitempty"`
}

type execResponse struct {
	Output    string         `json:"output"`
	ExitCode  int            `json:"exit_code"`
	Action    command.Action `json:"action,omitempty"`
	Verbose   string         `json:"verbose,omitempty"`
	SessionID string         `json:"session_id"`
}

type completeRequest struct {
	Line      string `json:"line"`
	Cursor    int    `json:"cursor"`
	SessionID string `json:"session_id,omitempty"`
}

type completeResponse struct {
	Start       int      `json:"start"`
	End         int      `json:"end"`
	Replacement string   `json:"replacement"`
	Candidates  []string `json:"candidates"`
	SessionID   string   `json:"session_id"`
}

func NewHandler(sessions *shell.SessionManager) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", health)
	mux.Handle(
		"GET /testui/",
		http.StripPrefix("/testui/", http.FileServer(http.Dir("testui"))),
	)
	mux.HandleFunc("POST /api/exec", func(w http.ResponseWriter, r *http.Request) {
		exec(w, r, sessions)
	})
	mux.HandleFunc("POST /api/complete", func(w http.ResponseWriter, r *http.Request) {
		complete(w, r, sessions)
	})

	return mux
}

func complete(w http.ResponseWriter, r *http.Request, sessions *shell.SessionManager) {
	var request completeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	result, sessionID, err := sessions.Complete(
		request.SessionID,
		request.Line,
		request.Cursor,
	)
	if err != nil {
		http.Error(w, "session creation failed", http.StatusInternalServerError)
		return
	}

	response := completeResponse{
		Start:       result.Start,
		End:         result.End,
		Replacement: result.Replacement,
		Candidates:  result.Candidates,
		SessionID:   sessionID,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

func health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("ok"))
}

func exec(w http.ResponseWriter, r *http.Request, sessions *shell.SessionManager) {
	var request execRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	result, sessionID, err := sessions.Execute(request.SessionID, request.Line)
	if err != nil {
		http.Error(w, "session creation failed", http.StatusInternalServerError)
		return
	}
	response := execResponse{
		Output:    result.Output,
		ExitCode:  result.ExitCode,
		Action:    result.Action,
		Verbose:   result.Verbose,
		SessionID: sessionID,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}
