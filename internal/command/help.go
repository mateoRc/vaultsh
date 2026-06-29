package command

import (
	"fmt"
	"strings"
)

type Help struct {
	registry *Registry
}

func NewHelp(registry *Registry) Help {
	return Help{registry: registry}
}

func (Help) Name() string {
	return "help"
}

func (Help) Description() string {
	return "List available commands"
}

func (h Help) Execute() Result {
	var output strings.Builder
	output.WriteString("Available commands:")
	for _, command := range h.registry.Commands() {
		fmt.Fprintf(&output, "\n  %s - %s", command.Name(), command.Description())
	}

	return Result{
		Output:   output.String(),
		ExitCode: 0,
	}
}
