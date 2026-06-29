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

func (Help) Usage() string {
	return "help [command]"
}

func (h Help) Execute(args []string, _ Input) Result {
	if len(args) > 1 {
		return Result{
			Output:   "usage: help [command]",
			ExitCode: ExitUsage,
		}
	}

	if len(args) == 1 {
		current, found := h.registry.Find(args[0])
		if !found {
			return Result{
				Output:   fmt.Sprintf("help: no help topic for %s", args[0]),
				ExitCode: ExitFailure,
			}
		}

		usage := current.Name()
		if provider, ok := current.(interface{ Usage() string }); ok {
			usage = provider.Usage()
		}

		return Result{
			Output:   fmt.Sprintf("Usage: %s\n%s", usage, current.Description()),
			ExitCode: ExitSuccess,
		}
	}

	var output strings.Builder
	output.WriteString("Available commands:")
	for _, command := range h.registry.Commands() {
		if IsHidden(command) {
			continue
		}
		fmt.Fprintf(&output, "\n  %s - %s", command.Name(), command.Description())
	}

	return Result{
		Output:   output.String(),
		ExitCode: ExitSuccess,
	}
}
