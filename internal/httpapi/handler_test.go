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
	oracle bool
}

type visitorRecorderStub struct {
	visits []string
}

func (r *visitorRecorderStub) Availability() (bool, bool, bool) { return false, true, false }
func (r *visitorRecorderStub) RecordVisit(visitorID string) error {
	r.visits = append(r.visits, visitorID)
	return nil
}

func TestVisitValidatesAndForwardsVisitorID(t *testing.T) {
	recorder := &visitorRecorderStub{}
	handler := NewHandlerWithStatus(shell.NewSessionManager(filesystem.NewDirectory("")), recorder)
	for _, test := range []struct {
		body string
		code int
	}{
		{`{"visitor_id":"invalid"}`, http.StatusBadRequest},
		{`{"visitor_id":"12345678-1234-4123-8123-123456789abc"}`, http.StatusNoContent},
	} {
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/api/visit", strings.NewReader(test.body)))
		if response.Code != test.code {
			t.Errorf("body %s: status = %d, want %d", test.body, response.Code, test.code)
		}
	}
	if len(recorder.visits) != 1 || recorder.visits[0] != "12345678-1234-4123-8123-123456789abc" {
		t.Errorf("recorded visits = %v", recorder.visits)
	}
}

const aboutOutput = "Vaultsh is the read-only virtual shell and Explorer for Backend Lab.\n\n" +
	"Start with /docs/projects/overview.md for the integrated system and " +
	"its boundaries. Professional background is under /cv."

func (s statusStub) Availability() (bool, bool, bool) {
	return s.atlas, s.forge, s.oracle
}

func TestStatusReportsExternalServices(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/status", nil)
	response := httptest.NewRecorder()
	NewHandlerWithStatus(
		shell.NewSessionManager(filesystem.NewDirectory("")),
		statusStub{atlas: true, forge: false, oracle: true},
	).ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d", response.Code)
	}
	if response.Body.String() != "{\"atlas\":true,\"forge\":false,\"oracle\":true}\n" {
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

func TestBrowseVirtualFiles(t *testing.T) {
	root := filesystem.NewDirectory("")
	cv := filesystem.NewDirectory("cv")
	if err := cv.Add(filesystem.NewFile("about.md", "# About\n\nA profile.")); err != nil {
		t.Fatal(err)
	}
	if err := cv.Add(filesystem.NewFile(".private.md", "hidden")); err != nil {
		t.Fatal(err)
	}
	if err := root.Add(cv); err != nil {
		t.Fatal(err)
	}
	handler := NewHandler(shell.NewSessionManager(root))

	for _, test := range []struct {
		path string
		code int
		kind filesystem.Kind
		children int
		content string
	}{
		{"/api/browse?path=/cv", http.StatusOK, filesystem.KindDirectory, 1, ""},
		{"/api/browse?path=/cv/about.md", http.StatusOK, filesystem.KindFile, 0, "# About\n\nA profile."},
		{"/api/browse?path=../../", http.StatusBadRequest, "", 0, ""},
		{"/api/browse?path=/missing", http.StatusNotFound, "", 0, ""},
	} {
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, test.path, nil))
		if response.Code != test.code {
			t.Errorf("%s: status = %d, want %d", test.path, response.Code, test.code)
			continue
		}
		if test.code != http.StatusOK {
			continue
		}
		var result shell.BrowseResult
		if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
			t.Fatal(err)
		}
		if result.Kind != test.kind || len(result.Children) != test.children || result.Content != test.content {
			t.Errorf("%s: unexpected result: %+v", test.path, result)
		}
		if test.children == 1 && len(result.Children) == 1 && result.Children[0].Path != "/cv/about.md" {
			t.Errorf("child path = %q", result.Children[0].Path)
		}
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
