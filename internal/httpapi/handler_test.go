package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mateom/vaultsh/internal/shell"
)

func TestHealth(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	response := httptest.NewRecorder()

	NewHandler(shell.New()).ServeHTTP(response, request)

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

	NewHandler(shell.New()).ServeHTTP(response, request)

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
}

func TestExecReturnsClearAction(t *testing.T) {
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/exec",
		strings.NewReader(`{"line":"clear"}`),
	)
	response := httptest.NewRecorder()

	NewHandler(shell.New()).ServeHTTP(response, request)

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

func TestExecRejectsInvalidJSON(t *testing.T) {
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/exec",
		strings.NewReader(`{"line":`),
	)
	response := httptest.NewRecorder()

	NewHandler(shell.New()).ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", response.Code, http.StatusBadRequest)
	}
}
