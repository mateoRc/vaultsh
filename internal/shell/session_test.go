package shell

import (
	"testing"

	"github.com/mateom/vaultsh/internal/filesystem"
)

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
