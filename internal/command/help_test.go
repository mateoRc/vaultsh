package command

import (
	"testing"

	"github.com/mateom/vaultsh/internal/filesystem"
)

func TestHelpIncludesUsageAndExample(t *testing.T) {
	registry := NewRegistry()
	registry.Register(NewCat(filesystem.NewWorkingDirectory(
		filesystem.NewDirectory(""),
	)))

	result := NewHelp(registry).Execute([]string{"cat"}, Input{})

	want := "Usage: cat [-n] [file]\n" +
		"Print file contents\n" +
		"Example: cat /cv/about.md"
	if result.ExitCode != ExitSuccess || result.Output != want {
		t.Errorf("result = %#v", result)
	}
}
