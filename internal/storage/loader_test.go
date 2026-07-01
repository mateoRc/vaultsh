package storage

import (
	"testing"
	"testing/fstest"

	"github.com/mateom/vaultsh/internal/filesystem"
)

func TestLoad(t *testing.T) {
	source := fstest.MapFS{
		"about.txt":            {Data: []byte("about")},
		"projects/vaultsh.txt": {Data: []byte("vaultsh")},
	}

	root, err := Load(source)
	if err != nil {
		t.Fatalf("Load(): %v", err)
	}

	resolver := filesystem.NewResolver(root)
	node, _, err := resolver.Resolve("/", "/projects/vaultsh.txt")
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
