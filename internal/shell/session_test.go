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
