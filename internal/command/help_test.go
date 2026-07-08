package command

import (
	"testing"

	"github.com/mateom/vaultsh/internal/filesystem"
)

type helpRegistryStub struct {
	commands []Command
}

func (r helpRegistryStub) Find(name string) (Command, bool) {
	for _, current := range r.commands {
		if current.Name() == name {
			return current, true
		}
	}
	return nil, false
}

func (r helpRegistryStub) Commands() []Command {
	return r.commands
}

func TestHelpIncludesUsageAndExample(t *testing.T) {
	registry := helpRegistryStub{commands: []Command{
		NewCat(filesystem.NewWorkingDirectory(filesystem.NewDirectory(""))),
	}}

	result := NewHelp(registry).Execute([]string{"cat"}, Input{})

	want := "Usage: cat [-n] [file]\n" +
		"Print file contents\n" +
		"Example: cat /cv/about.md"
	if result.ExitCode != ExitSuccess || result.Output != want {
		t.Errorf("result = %#v", result)
	}
}
