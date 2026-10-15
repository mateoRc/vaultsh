package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/mateom/vaultsh/internal/command"
	"github.com/mateom/vaultsh/internal/shell"
)

const (
	maxRequestBodyBytes = 16 * 1024
	maxCommandLength    = 4096
	jsonContentType     = "application/json"
)

type routeLimit struct {
	keyPrefix string
	perMinute int
	burst     int
	retry     string
}

var (
	statusLimit = routeLimit{
		keyPrefix: "status",
		perMinute: 120,
		burst:     20,
		retry:     "1",
	}
	execLimit = routeLimit{
		keyPrefix: "exec",
		perMinute: 30,
		burst:     10,
		retry:     "2",
	}
	completeLimit = routeLimit{
		keyPrefix: "complete",
		perMinute: 120,
		burst:     20,
		retry:     "1",
	}
	assistantLimit = routeLimit{
		keyPrefix: "assistant",
		perMinute: 10,
		burst:     3,
		retry:     "10",
	}
)

type HandlerConfig struct {
	TrustProxyHeaders bool
}

type execRequest struct {
	Line             string `json:"line"`
	SessionID        string `json:"session_id,omitempty"`
	AutomaticWelcome bool   `json:"automatic_welcome,omitempty"`
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

type VisitorRecorder interface {
	RecordVisit(visitorID string) error
}

var visitorIDPattern = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

type visitRequest struct {
	VisitorID string `json:"visitor_id"`
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
	limited := func(limit routeLimit, next http.HandlerFunc) http.HandlerFunc {
		return withRateLimit(limiter, config.TrustProxyHeaders, limit, next)
	}

	mux.HandleFunc("GET /", redirectToVault)
	mux.HandleFunc("GET /healthz", health)
	mux.HandleFunc("GET /api/status", limited(statusLimit, statusHandler(status)))
	mux.HandleFunc("POST /api/visit", limited(statusLimit, func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)
		var request visitRequest
		if !decodeRequest(w, r, &request) {
			return
		}
		if !visitorIDPattern.MatchString(request.VisitorID) {
			http.Error(w, "invalid visitor ID", http.StatusBadRequest)
			return
		}
		recorder, ok := status.(VisitorRecorder)
		if !ok || recorder.RecordVisit(request.VisitorID) != nil {
			http.Error(w, "analytics unavailable", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	mux.HandleFunc("GET /api/browse", limited(statusLimit, func(w http.ResponseWriter, r *http.Request) {
		browse(w, r, sessions)
	}))
	mux.Handle(
		"GET /vault/",
		http.StripPrefix("/vault/", http.FileServer(http.Dir("web"))),
	)
	mux.HandleFunc("POST /api/exec", limited(execLimit, func(w http.ResponseWriter, r *http.Request) {
		exec(w, r, sessions, limiter, config.TrustProxyHeaders)
	}))
	mux.HandleFunc("POST /api/complete", limited(completeLimit, func(w http.ResponseWriter, r *http.Request) {
		complete(w, r, sessions)
	}))

	return securityHeaders(mux)
}

func redirectToVault(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, "/vault/", http.StatusTemporaryRedirect)
}

func statusHandler(status StatusProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		atlas, forge := false, false
		if status != nil {
			atlas, forge = status.Availability()
		}
		writeJSON(w, map[string]bool{"atlas": atlas, "forge": forge})
	}
}

func browse(w http.ResponseWriter, r *http.Request, sessions *shell.SessionManager) {
	target := r.URL.Query().Get("path")
	if target == "" {
		target = "/"
	}
	if len(target) > maxCommandLength || !strings.HasPrefix(target, "/") {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	result, err := sessions.Browse(target)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	writeJSON(w, result)
}

func withRateLimit(
	limiter *rateLimiter,
	trustProxyHeaders bool,
	limit routeLimit,
	next http.HandlerFunc,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := limit.keyPrefix + ":" + clientIP(r, trustProxyHeaders)
		if limiter.allow(key, limit.perMinute, limit.burst) {
			next(w, r)
			return
		}
		w.Header().Set("Retry-After", limit.retry)
		http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
	}
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

	writeJSON(w, response)
}

func health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("ok"))
}

func exec(
	w http.ResponseWriter,
	r *http.Request,
	sessions *shell.SessionManager,
	limiter *rateLimiter,
	trustProxyHeaders bool,
) {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)
	var request execRequest
	if !decodeRequest(w, r, &request) {
		return
	}
	if len(request.Line) > maxCommandLength {
		http.Error(w, "command is too long", http.StatusRequestEntityTooLarge)
		return
	}
	if isAssistantCommand(request.Line) &&
		!limiter.allow("assistant:"+clientIP(r, trustProxyHeaders), assistantLimit.perMinute, assistantLimit.burst) {
		w.Header().Set("Retry-After", assistantLimit.retry)
		http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	result, sessionID, currentDirectory, err := sessions.ExecuteWithOptions(
		request.SessionID,
		request.Line,
		request.AutomaticWelcome,
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

	writeJSON(w, response)
}

func isAssistantCommand(line string) bool {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return false
	}
	return fields[0] == "invoke" || fields[0] == "ask"
}

func writeJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", jsonContentType)
	_ = json.NewEncoder(w).Encode(value)
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
		if strings.HasPrefix(r.URL.Path, "/vault/") {
			w.Header().Set("Cache-Control", "no-cache")
		}
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
