package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/mateom/vaultsh/internal/command"
	"github.com/mateom/vaultsh/internal/shell"
)

const (
	maxRequestBodyBytes = 16 * 1024
	maxCommandLength    = 4096
)

type HandlerConfig struct {
	TrustProxyHeaders bool
}

type execRequest struct {
	Line      string `json:"line"`
	SessionID string `json:"session_id,omitempty"`
}

type execResponse struct {
	Output           string         `json:"output"`
	ExitCode         int            `json:"exit_code"`
	Action           command.Action `json:"action,omitempty"`
	Verbose          string         `json:"verbose,omitempty"`
	SessionID        string         `json:"session_id"`
	CurrentDirectory string         `json:"current_directory"`
}

type completeRequest struct {
	Line      string `json:"line"`
	Cursor    int    `json:"cursor"`
	SessionID string `json:"session_id,omitempty"`
}

type completeResponse struct {
	Start            int      `json:"start"`
	End              int      `json:"end"`
	Replacement      string   `json:"replacement"`
	Candidates       []string `json:"candidates"`
	SessionID        string   `json:"session_id"`
	CurrentDirectory string   `json:"current_directory"`
}

type StatusProvider interface {
	Availability() (atlas bool, forge bool)
}

func NewHandler(sessions *shell.SessionManager) http.Handler {
	return NewHandlerWithStatus(sessions, nil)
}

func NewHandlerWithStatus(
	sessions *shell.SessionManager,
	status StatusProvider,
) http.Handler {
	return NewHandlerWithConfig(sessions, status, HandlerConfig{})
}

func NewHandlerWithConfig(
	sessions *shell.SessionManager,
	status StatusProvider,
	config HandlerConfig,
) http.Handler {
	mux := http.NewServeMux()
	limiter := newRateLimiter()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, "/vault/", http.StatusTemporaryRedirect)
	})
	mux.HandleFunc("GET /healthz", health)
	mux.HandleFunc("GET /api/status", func(w http.ResponseWriter, r *http.Request) {
		if !limiter.allow("status:"+clientIP(r, config.TrustProxyHeaders), 120, 20) {
			w.Header().Set("Retry-After", "1")
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		atlas, forge := false, false
		if status != nil {
			atlas, forge = status.Availability()
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]bool{
			"atlas": atlas,
			"forge": forge,
		})
	})
	mux.Handle(
		"GET /vault/",
		http.StripPrefix("/vault/", http.FileServer(http.Dir("web"))),
	)
	mux.HandleFunc("POST /api/exec", func(w http.ResponseWriter, r *http.Request) {
		if !limiter.allow("exec:"+clientIP(r, config.TrustProxyHeaders), 30, 10) {
			w.Header().Set("Retry-After", "2")
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		exec(w, r, sessions)
	})
	mux.HandleFunc("POST /api/complete", func(w http.ResponseWriter, r *http.Request) {
		if !limiter.allow("complete:"+clientIP(r, config.TrustProxyHeaders), 120, 20) {
			w.Header().Set("Retry-After", "1")
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		complete(w, r, sessions)
	})

	return securityHeaders(mux)
}

func complete(w http.ResponseWriter, r *http.Request, sessions *shell.SessionManager) {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)
	var request completeRequest
	if !decodeRequest(w, r, &request) {
		return
	}
	if len(request.Line) > maxCommandLength {
		http.Error(w, "input is too long", http.StatusRequestEntityTooLarge)
		return
	}

	result, sessionID, currentDirectory, err := sessions.Complete(
		request.SessionID,
		request.Line,
		request.Cursor,
	)
	if err != nil {
		writeSessionError(w, err)
		return
	}

	response := completeResponse{
		Start:            result.Start,
		End:              result.End,
		Replacement:      result.Replacement,
		Candidates:       result.Candidates,
		SessionID:        sessionID,
		CurrentDirectory: currentDirectory,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

func health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("ok"))
}

func exec(w http.ResponseWriter, r *http.Request, sessions *shell.SessionManager) {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)
	var request execRequest
	if !decodeRequest(w, r, &request) {
		return
	}
	if len(request.Line) > maxCommandLength {
		http.Error(w, "command is too long", http.StatusRequestEntityTooLarge)
		return
	}

	result, sessionID, currentDirectory, err := sessions.Execute(
		request.SessionID,
		request.Line,
	)
	if err != nil {
		writeSessionError(w, err)
		return
	}
	response := execResponse{
		Output:           result.Output,
		ExitCode:         result.ExitCode,
		Action:           result.Action,
		Verbose:          result.Verbose,
		SessionID:        sessionID,
		CurrentDirectory: currentDirectory,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

func decodeRequest(w http.ResponseWriter, r *http.Request, target any) bool {
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			http.Error(w, "request body is too large", http.StatusRequestEntityTooLarge)
			return false
		}
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return false
	}
	return true
}

func writeSessionError(w http.ResponseWriter, err error) {
	if errors.Is(err, shell.ErrSessionLimit) {
		w.Header().Set("Retry-After", "60")
		http.Error(w, "session capacity reached", http.StatusTooManyRequests)
		return
	}
	http.Error(w, "session creation failed", http.StatusInternalServerError)
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(
			"Content-Security-Policy",
			"default-src 'self'; base-uri 'none'; frame-ancestors 'none'; "+
				"form-action 'self'; object-src 'none'",
		)
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		next.ServeHTTP(w, r)
	})
}
