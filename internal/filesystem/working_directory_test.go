package filesystem

import (
	"errors"
	"testing"
)

func TestWorkingDirectoryStartsAtRoot(t *testing.T) {
	root := NewDirectory("")
	workingDirectory := NewWorkingDirectory(root)

	if workingDirectory.Path() != "/" {
		t.Errorf("Path() = %q, want %q", workingDirectory.Path(), "/")
	}
	if workingDirectory.Directory() != root {
		t.Errorf("Directory() = %#v, want root", workingDirectory.Directory())
	}
}

func TestWorkingDirectoryChange(t *testing.T) {
	root := NewDirectory("")
	docs := NewDirectory("docs")
	mustAdd(t, root, docs)
	workingDirectory := NewWorkingDirectory(root)

	if err := workingDirectory.Change("docs"); err != nil {
		t.Fatalf("Change(docs): %v", err)
	}

	if workingDirectory.Path() != "/docs" {
		t.Errorf("Path() = %q, want %q", workingDirectory.Path(), "/docs")
	}
	if workingDirectory.Directory() != docs {
		t.Errorf("Directory() = %#v, want docs", workingDirectory.Directory())
	}
}

func TestWorkingDirectoryResolvesRelativePath(t *testing.T) {
	root := NewDirectory("")
	docs := NewDirectory("docs")
	readme := NewFile("readme.txt", "hello")
	mustAdd(t, root, docs)
	mustAdd(t, docs, readme)
	workingDirectory := NewWorkingDirectory(root)
	if err := workingDirectory.Change("/docs"); err != nil {
		t.Fatalf("Change(/docs): %v", err)
	}

	node, resolvedPath, err := workingDirectory.Resolve("readme.txt")

	if err != nil {
		t.Fatalf("Resolve(readme.txt): %v", err)
	}
	if node != readme {
		t.Errorf("node = %#v, want readme", node)
	}
	if resolvedPath != "/docs/readme.txt" {
		t.Errorf("path = %q, want %q", resolvedPath, "/docs/readme.txt")
	}
}

func TestWorkingDirectoryRejectsFile(t *testing.T) {
	root := NewDirectory("")
	mustAdd(t, root, NewFile("readme.txt", "hello"))
	workingDirectory := NewWorkingDirectory(root)

	err := workingDirectory.Change("/readme.txt")

	if !errors.Is(err, ErrNotDirectory) {
		t.Errorf("Change() error = %v, want %v", err, ErrNotDirectory)
	}
	if workingDirectory.Path() != "/" {
		t.Errorf("Path() = %q after failed change, want /", workingDirectory.Path())
	}
}

func TestWorkingDirectoryPreservesPathOnMissingTarget(t *testing.T) {
	workingDirectory := NewWorkingDirectory(NewDirectory(""))

	err := workingDirectory.Change("/missing")

	if !errors.Is(err, ErrNodeNotFound) {
		t.Errorf("Change() error = %v, want %v", err, ErrNodeNotFound)
	}
	if workingDirectory.Path() != "/" {
		t.Errorf("Path() = %q after failed change, want /", workingDirectory.Path())
	}
}
