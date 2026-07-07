package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mateom/vaultsh/internal/filesystem"
	"github.com/mateom/vaultsh/internal/shell"
)

type statusStub struct {
	atlas bool
	forge bool
}

const aboutOutput = "Vaultsh is a read-only virtual shell for exploring Mateo's " +
	"CV, projects, and live backend services.\n\n" +
	"The portfolio is the system: a Go shell backed by " +
	"Atlas search, Forge telemetry, and Lab deployment docs."

func (s statusStub) Availability() (bool, bool) {
	return s.atlas, s.forge
}

func TestStatusReportsExternalServices(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/status", nil)
	response := httptest.NewRecorder()
	NewHandlerWithStatus(
		shell.NewSessionManager(filesystem.NewDirectory("")),
		statusStub{atlas: true, forge: false},
	).ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d", response.Code)
	}
	if response.Body.String() != "{\"atlas\":true,\"forge\":false}\n" {
		t.Errorf("body = %q", response.Body.String())
	}
}

func TestHealth(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	response := httptest.NewRecorder()

	newTestHandler().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", response.Code, http.StatusOK)
	}
	if response.Body.String() != "ok" {
		t.Errorf("body = %q, want %q", response.Body.String(), "ok")
	}
}

func TestRootRedirectsToVault(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	newTestHandler().ServeHTTP(response, request)

	if response.Code != http.StatusTemporaryRedirect {
		t.Errorf("status = %d, want %d", response.Code, http.StatusTemporaryRedirect)
	}
	if location := response.Header().Get("Location"); location != "/vault/" {
		t.Errorf("Location = %q, want %q", location, "/vault/")
	}
}

func TestExec(t *testing.T) {
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/exec",
		strings.NewReader(`{"line":"about"}`),
	)
	response := httptest.NewRecorder()

	newTestHandler().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}

	var result execResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if result.Output != aboutOutput {
		t.Errorf("output = %q, want about output", result.Output)
	}
	if result.ExitCode != 0 {
		t.Errorf("exit code = %d, want 0", result.ExitCode)
	}
	if result.SessionID == "" {
		t.Error("session ID is empty")
	}
	if result.CurrentDirectory != "/" {
		t.Errorf("current directory = %q, want /", result.CurrentDirectory)
	}
}

func TestExecReturnsUpdatedCurrentDirectory(t *testing.T) {
	root := filesystem.NewDirectory("")
	if err := root.Add(filesystem.NewDirectory("docs")); err != nil {
		t.Fatal(err)
	}
	handler := NewHandler(shell.NewSessionManager(root))
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/exec",
		strings.NewReader(`{"line":"cd /docs"}`),
	)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	var result execResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if result.CurrentDirectory != "/docs" {
		t.Errorf("current directory = %q, want /docs", result.CurrentDirectory)
	}
}

func TestExecReturnsClearAction(t *testing.T) {
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/exec",
		strings.NewReader(`{"line":"clear"}`),
	)
	response := httptest.NewRecorder()

	newTestHandler().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}

	var result execResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if result.Action != "clear" {
		t.Errorf("action = %q, want %q", result.Action, "clear")
	}
}

func TestExecReturnsVerboseDetailsOnlyWhenRequested(t *testing.T) {
	tests := []struct {
		name        string
		line        string
		wantVerbose string
	}{
		{
			name: "normal response",
			line: "pwd",
		},
		{
			name:        "verbose response",
			line:        "pwd --verbose",
			wantVerbose: "pipeline=pwd; stages=1; completed=1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(
				http.MethodPost,
				"/api/exec",
				strings.NewReader(`{"line":"`+tt.line+`"}`),
			)
			response := httptest.NewRecorder()

			newTestHandler().ServeHTTP(response, request)

			var result map[string]any
			if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
				t.Fatalf("decode response: %v", err)
			}
			got, present := result["verbose"]
			if tt.wantVerbose == "" {
				if present {
					t.Errorf("verbose field = %q, want omitted", got)
				}
				return
			}
			if got != tt.wantVerbose {
				t.Errorf("verbose field = %q, want %q", got, tt.wantVerbose)
			}
		})
	}
}

func TestExecRejectsInvalidJSON(t *testing.T) {
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/exec",
		strings.NewReader(`{"line":`),
	)
	response := httptest.NewRecorder()

	newTestHandler().ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", response.Code, http.StatusBadRequest)
	}
}

func TestComplete(t *testing.T) {
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/complete",
		strings.NewReader(`{"line":"ca","cursor":2}`),
	)
	response := httptest.NewRecorder()

	newTestHandler().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}

	var result completeResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if result.Replacement != "cat " {
		t.Errorf("replacement = %q, want %q", result.Replacement, "cat ")
	}
	if result.SessionID == "" {
		t.Error("session ID is empty")
	}
	if result.CurrentDirectory != "/" {
		t.Errorf("current directory = %q, want /", result.CurrentDirectory)
	}
}

func TestExecRejectsOversizedBody(t *testing.T) {
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/exec",
		strings.NewReader(`{"line":"`+strings.Repeat("x", maxRequestBodyBytes)+`"}`),
	)
	response := httptest.NewRecorder()

	newTestHandler().ServeHTTP(response, request)

	if response.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("status = %d, want %d", response.Code, http.StatusRequestEntityTooLarge)
	}
}

func TestExecRejectsLongCommand(t *testing.T) {
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/exec",
		strings.NewReader(`{"line":"`+strings.Repeat("x", maxCommandLength+1)+`"}`),
	)
	response := httptest.NewRecorder()

	newTestHandler().ServeHTTP(response, request)

	if response.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("status = %d, want %d", response.Code, http.StatusRequestEntityTooLarge)
	}
}

func TestExecRateLimit(t *testing.T) {
	handler := newTestHandler()
	for index := 0; index < 11; index++ {
		request := httptest.NewRequest(
			http.MethodPost,
			"/api/exec",
			strings.NewReader(`{"line":"pwd"}`),
		)
		request.RemoteAddr = "192.0.2.1:1234"
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)

		if index < 10 && response.Code != http.StatusOK {
			t.Fatalf("request %d status = %d", index+1, response.Code)
		}
		if index == 10 && response.Code != http.StatusTooManyRequests {
			t.Errorf("request 11 status = %d, want %d", response.Code, http.StatusTooManyRequests)
		}
	}
}

func TestSecurityHeaders(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	response := httptest.NewRecorder()

	newTestHandler().ServeHTTP(response, request)

	if response.Header().Get("Content-Security-Policy") == "" {
		t.Error("Content-Security-Policy header is missing")
	}
	if response.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Error("X-Content-Type-Options header is missing")
	}
}

func newTestHandler() http.Handler {
	root := filesystem.NewDirectory("")
	return NewHandler(shell.NewSessionManager(root))
}
