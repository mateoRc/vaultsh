package shell

import (
	"reflect"
	"testing"

	"github.com/mateom/vaultsh/internal/filesystem"
)

func TestCompleteCommand(t *testing.T) {
	engine := New()

	result := engine.Complete("ca", len("ca"))

	if result.Replacement != "cat " {
		t.Errorf("replacement = %q, want %q", result.Replacement, "cat ")
	}
	if !reflect.DeepEqual(result.Candidates, []string{"cat"}) {
		t.Errorf("candidates = %q, want [cat]", result.Candidates)
	}
}

func TestCompleteCommandCommonPrefix(t *testing.T) {
	engine := New()

	result := engine.Complete("c", len("c"))

	if result.Replacement != "c" {
		t.Errorf("replacement = %q, want %q", result.Replacement, "c")
	}
	if !reflect.DeepEqual(result.Candidates, []string{"cat", "cd", "clear"}) {
		t.Errorf("candidates = %q, want [cat cd clear]", result.Candidates)
	}
}

func TestCompletePath(t *testing.T) {
	root := filesystem.NewDirectory("")
	if err := root.Add(filesystem.NewDirectory("experience")); err != nil {
		t.Fatalf("Add(experience): %v", err)
	}
	engine := NewWithRoot(root)

	result := engine.Complete("cd exp", len("cd exp"))

	if result.Replacement != "experience/" {
		t.Errorf("replacement = %q, want %q", result.Replacement, "experience/")
	}
	if result.Start != len("cd ") || result.End != len("cd exp") {
		t.Errorf("replacement range = %d:%d", result.Start, result.End)
	}
}

func TestCompletePathUsesWorkingDirectory(t *testing.T) {
	root := filesystem.NewDirectory("")
	experience := filesystem.NewDirectory("experience")
	if err := root.Add(experience); err != nil {
		t.Fatalf("Add(experience): %v", err)
	}
	if err := experience.Add(filesystem.NewFile("reversinglabs.txt", "")); err != nil {
		t.Fatalf("Add(reversinglabs.txt): %v", err)
	}
	engine := NewWithRoot(root)
	if result := engine.Execute("cd experience"); result.ExitCode != 0 {
		t.Fatalf("cd failed: %s", result.Output)
	}

	result := engine.Complete("cat rev", len("cat rev"))

	if result.Replacement != "reversinglabs.txt " {
		t.Errorf(
			"replacement = %q, want %q",
			result.Replacement,
			"reversinglabs.txt ",
		)
	}
}
