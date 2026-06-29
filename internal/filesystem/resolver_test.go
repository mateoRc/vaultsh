package filesystem

import (
	"errors"
	"testing"
)

func TestResolverResolve(t *testing.T) {
	root := NewDirectory("")
	docs := NewDirectory("docs")
	readme := NewFile("readme.txt", "hello")
	mustAdd(t, root, docs)
	mustAdd(t, docs, readme)

	tests := []struct {
		name             string
		workingDirectory string
		target           string
		wantNode         Node
		wantPath         string
	}{
		{
			name:     "root",
			target:   "/",
			wantNode: root,
			wantPath: "/",
		},
		{
			name:     "absolute path",
			target:   "/docs/readme.txt",
			wantNode: readme,
			wantPath: "/docs/readme.txt",
		},
		{
			name:             "relative path",
			workingDirectory: "/docs",
			target:           "readme.txt",
			wantNode:         readme,
			wantPath:         "/docs/readme.txt",
		},
		{
			name:             "parent directory",
			workingDirectory: "/docs",
			target:           "..",
			wantNode:         root,
			wantPath:         "/",
		},
		{
			name:     "clean path",
			target:   "/docs//./readme.txt",
			wantNode: readme,
			wantPath: "/docs/readme.txt",
		},
	}

	resolver := NewResolver(root)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, resolvedPath, err := resolver.Resolve(tt.workingDirectory, tt.target)
			if err != nil {
				t.Fatalf("Resolve(): %v", err)
			}
			if node != tt.wantNode {
				t.Errorf("node = %#v, want %#v", node, tt.wantNode)
			}
			if resolvedPath != tt.wantPath {
				t.Errorf("path = %q, want %q", resolvedPath, tt.wantPath)
			}
		})
	}
}

func TestResolverMissingNode(t *testing.T) {
	resolver := NewResolver(NewDirectory(""))

	_, _, err := resolver.Resolve("/", "/missing")

	if !errors.Is(err, ErrNodeNotFound) {
		t.Errorf("Resolve() error = %v, want %v", err, ErrNodeNotFound)
	}
}

func TestResolverCannotTraverseFile(t *testing.T) {
	root := NewDirectory("")
	mustAdd(t, root, NewFile("readme.txt", "hello"))
	resolver := NewResolver(root)

	_, _, err := resolver.Resolve("/", "/readme.txt/child")

	if !errors.Is(err, ErrNotDirectory) {
		t.Errorf("Resolve() error = %v, want %v", err, ErrNotDirectory)
	}
}

func mustAdd(t *testing.T, directory *Directory, node Node) {
	t.Helper()
	if err := directory.Add(node); err != nil {
		t.Fatalf("Add(%s): %v", node.Name(), err)
	}
}
