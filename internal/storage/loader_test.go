package storage

import (
	"testing"
	"testing/fstest"
	"time"

	"github.com/mateom/vaultsh/internal/filesystem"
)

func TestLoad(t *testing.T) {
	source := fstest.MapFS{
		"about.md":            {Data: []byte("about")},
		"projects/vaultsh.md": {Data: []byte("vaultsh")},
	}

	root, err := Load(source)
	if err != nil {
		t.Fatalf("Load(): %v", err)
	}

	resolver := filesystem.NewResolver(root)
	node, _, err := resolver.Resolve("/", "/projects/vaultsh.md")
	if err != nil {
		t.Fatalf("Resolve(): %v", err)
	}

	file, ok := node.(*filesystem.File)
	if !ok {
		t.Fatalf("node type = %T, want *filesystem.File", node)
	}
	if file.Content() != "vaultsh" {
		t.Errorf("Content() = %q, want %q", file.Content(), "vaultsh")
	}
}

func TestLoadPreservesModTime(t *testing.T) {
	modTime := time.Date(2026, 7, 7, 12, 0, 0, 0, time.UTC)
	source := fstest.MapFS{
		"about.md": {Data: []byte("about"), ModTime: modTime},
	}

	root, err := Load(source)
	if err != nil {
		t.Fatalf("Load(): %v", err)
	}

	resolver := filesystem.NewResolver(root)
	node, _, err := resolver.Resolve("/", "/about.md")
	if err != nil {
		t.Fatalf("Resolve(): %v", err)
	}
	if !node.ModTime().Equal(modTime) {
		t.Errorf("ModTime() = %s, want %s", node.ModTime(), modTime)
	}
}
