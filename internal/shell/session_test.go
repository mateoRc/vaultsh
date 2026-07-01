package shell

import (
	"context"
	"testing"
	"time"

	"github.com/mateom/vaultsh/internal/command"
	"github.com/mateom/vaultsh/internal/filesystem"
)

type eventRecorderStub struct {
	service  string
	event    string
	name     string
	exitCode int
}

func (r *eventRecorderStub) Record(
	service string,
	event string,
	name string,
	_ int64,
	exitCode int,
) error {
	r.service = service
	r.event = event
	r.name = name
	r.exitCode = exitCode
	return nil
}

type unavailableServices struct{}

func (unavailableServices) Search(string) ([]command.SearchResult, error) {
	return nil, context.DeadlineExceeded
}

func (unavailableServices) Summary() (command.MetricsSummary, error) {
	return command.MetricsSummary{}, context.DeadlineExceeded
}

func (unavailableServices) Dashboard() (string, error) {
	return "", context.DeadlineExceeded
}

func TestSessionManagerRecordsCommandTelemetryWithoutChangingResult(t *testing.T) {
	recorder := &eventRecorderStub{}
	manager := NewSessionManagerWithDependencies(
		filesystem.NewDirectory(""),
		Dependencies{Events: recorder},
	)

	result, _, err := manager.Execute("", "about")

	if err != nil || result.ExitCode != command.ExitSuccess {
		t.Fatalf("Execute() result = %#v, error = %v", result, err)
	}
	if recorder.service != "vault" ||
		recorder.event != "command.executed" ||
		recorder.name != "about" ||
		recorder.exitCode != command.ExitSuccess {
		t.Errorf("recorded event = %#v", recorder)
	}
}

func TestUnavailableIntegrationsDoNotBreakCoreCommands(t *testing.T) {
	services := unavailableServices{}
	manager := NewSessionManagerWithDependencies(
		filesystem.NewDirectory(""),
		Dependencies{Search: services, Metrics: services},
	)

	core, sessionID, err := manager.Execute("", "about")
	if err != nil || core.ExitCode != command.ExitSuccess {
		t.Fatalf("core command result = %#v, error = %v", core, err)
	}

	search, _, err := manager.Execute(sessionID, "search kafka")
	if err != nil || search.Output != "search unavailable" {
		t.Fatalf("search result = %#v, error = %v", search, err)
	}
}

func TestSessionManagerKeepsWorkingDirectoriesIndependent(t *testing.T) {
	root := filesystem.NewDirectory("")
	if err := root.Add(filesystem.NewDirectory("docs")); err != nil {
		t.Fatalf("Add(docs): %v", err)
	}
	manager := NewSessionManager(root)

	result, firstID, err := manager.Execute("", "cd docs")
	if err != nil {
		t.Fatalf("first Execute(): %v", err)
	}
	if result.ExitCode != 0 {
		t.Fatalf("cd exit code = %d, output = %q", result.ExitCode, result.Output)
	}

	result, returnedID, err := manager.Execute(firstID, "pwd")
	if err != nil {
		t.Fatalf("first session pwd: %v", err)
	}
	if returnedID != firstID {
		t.Errorf("returned session ID = %q, want %q", returnedID, firstID)
	}
	if result.Output != "/docs" {
		t.Errorf("first session pwd = %q, want /docs", result.Output)
	}

	result, secondID, err := manager.Execute("", "pwd")
	if err != nil {
		t.Fatalf("second session pwd: %v", err)
	}
	if secondID == "" || secondID == firstID {
		t.Errorf("second session ID = %q, want a new ID", secondID)
	}
	if result.Output != "/" {
		t.Errorf("second session pwd = %q, want /", result.Output)
	}
}

func TestSessionManagerKeepsHistoriesIndependent(t *testing.T) {
	manager := NewSessionManager(filesystem.NewDirectory(""))

	_, firstID, err := manager.Execute("", "pwd")
	if err != nil {
		t.Fatalf("first Execute(): %v", err)
	}
	_, secondID, err := manager.Execute("", "about")
	if err != nil {
		t.Fatalf("second Execute(): %v", err)
	}

	firstHistory, _, err := manager.Execute(firstID, "history")
	if err != nil {
		t.Fatalf("first history: %v", err)
	}
	secondHistory, _, err := manager.Execute(secondID, "history")
	if err != nil {
		t.Fatalf("second history: %v", err)
	}

	if firstHistory.Output != "1  pwd\n2  history" {
		t.Errorf("first history = %q", firstHistory.Output)
	}
	if secondHistory.Output != "1  about\n2  history" {
		t.Errorf("second history = %q", secondHistory.Output)
	}
}

func TestSessionManagerReplacesUnknownSessionID(t *testing.T) {
	manager := NewSessionManager(filesystem.NewDirectory(""))

	_, sessionID, err := manager.Execute("unknown", "pwd")

	if err != nil {
		t.Fatalf("Execute(): %v", err)
	}
	if sessionID == "" || sessionID == "unknown" {
		t.Errorf("session ID = %q, want a new server-generated ID", sessionID)
	}
}

func TestSessionManagerCleansUpExpiredSessions(t *testing.T) {
	now := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	manager := NewSessionManagerWithConfig(
		filesystem.NewDirectory(""),
		SessionConfig{
			TTL: time.Hour,
			Now: func() time.Time {
				return now
			},
		},
	)

	_, sessionID, err := manager.Execute("", "pwd")
	if err != nil {
		t.Fatalf("Execute(): %v", err)
	}

	now = now.Add(time.Hour)
	if removed := manager.CleanupExpired(); removed != 1 {
		t.Errorf("CleanupExpired() = %d, want 1", removed)
	}

	_, replacementID, err := manager.Execute(sessionID, "pwd")
	if err != nil {
		t.Fatalf("Execute(expired): %v", err)
	}
	if replacementID == sessionID {
		t.Error("expired session ID was reused")
	}
}

func TestSessionManagerRefreshesActivity(t *testing.T) {
	now := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	manager := NewSessionManagerWithConfig(
		filesystem.NewDirectory(""),
		SessionConfig{
			TTL: time.Hour,
			Now: func() time.Time {
				return now
			},
		},
	)

	_, sessionID, err := manager.Execute("", "pwd")
	if err != nil {
		t.Fatalf("first Execute(): %v", err)
	}
	now = now.Add(50 * time.Minute)
	if _, _, err := manager.Complete(sessionID, "c", 1); err != nil {
		t.Fatalf("Complete(): %v", err)
	}
	now = now.Add(20 * time.Minute)

	if removed := manager.CleanupExpired(); removed != 0 {
		t.Errorf("CleanupExpired() = %d, want 0", removed)
	}
}

func TestSessionManagerRejectsExpiredSessionBeforeCleanup(t *testing.T) {
	now := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	manager := NewSessionManagerWithConfig(
		filesystem.NewDirectory(""),
		SessionConfig{
			TTL: time.Hour,
			Now: func() time.Time {
				return now
			},
		},
	)

	_, sessionID, err := manager.Execute("", "pwd")
	if err != nil {
		t.Fatalf("first Execute(): %v", err)
	}
	now = now.Add(time.Hour)

	_, replacementID, err := manager.Execute(sessionID, "pwd")
	if err != nil {
		t.Fatalf("second Execute(): %v", err)
	}
	if replacementID == sessionID {
		t.Error("expired session ID was reused before cleanup")
	}
}

func TestSessionManagerCleanupStopsWithContext(t *testing.T) {
	manager := NewSessionManager(filesystem.NewDirectory(""))
	ctx, cancel := context.WithCancel(context.Background())
	stopped := make(chan struct{})

	go func() {
		manager.RunCleanup(ctx, time.Hour)
		close(stopped)
	}()
	cancel()

	select {
	case <-stopped:
	case <-time.After(time.Second):
		t.Fatal("RunCleanup() did not stop after context cancellation")
	}
}
