package sqlite_test

import (
	"context"
	"testing"

	embeddedcontent "github.com/mateom/vaultsh/content"
	contentsqlite "github.com/mateom/vaultsh/internal/content/sqlite"
)

func TestProviderLoadsEmbeddedCatalog(t *testing.T) {
	provider, err := contentsqlite.Open(embeddedcontent.Database)
	if err != nil {
		t.Fatalf("Open(): %v", err)
	}
	t.Cleanup(func() {
		if err := provider.Close(); err != nil {
			t.Errorf("Close(): %v", err)
		}
	})

	catalog, err := provider.Load(context.Background())
	if err != nil {
		t.Fatalf("Load(): %v", err)
	}

	if catalog.About.Text == "" {
		t.Error("About.Text is empty")
	}
	if len(catalog.Experiences) != 4 {
		t.Errorf("len(Experiences) = %d, want 4", len(catalog.Experiences))
	}
	if len(catalog.Projects) != 1 {
		t.Errorf("len(Projects) = %d, want 1", len(catalog.Projects))
	}
	if catalog.Skills.Text == "" {
		t.Error("Skills.Text is empty")
	}
}
