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
	if result.Output != "Vaultsh is a read-only virtual shell engine." {
		t.Errorf("output = %q, want about output", result.Output)
	}
	if result.ExitCode != 0 {
		t.Errorf("exit code = %d, want 0", result.ExitCode)
	}
	if result.SessionID == "" {
		t.Error("session ID is empty")
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
}

func newTestHandler() http.Handler {
	root := filesystem.NewDirectory("")
	return NewHandler(shell.NewSessionManager(root))
}
