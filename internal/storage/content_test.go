package storage

import (
	"context"
	"testing"

	"github.com/mateom/vaultsh/internal/content"
	"github.com/mateom/vaultsh/internal/filesystem"
)

type stubProvider struct {
	catalog content.Catalog
}

func (p stubProvider) Load(context.Context) (content.Catalog, error) {
	return p.catalog, nil
}

func TestLoad(t *testing.T) {
	provider := stubProvider{catalog: content.Catalog{
		About: content.About{Text: "about"},
		Experiences: []content.Experience{
			{Slug: "example", Text: "experience"},
		},
		Projects: []content.Project{
			{Slug: "vaultsh", Text: "vaultsh"},
		},
	}}

	root, err := Load(context.Background(), provider)
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
