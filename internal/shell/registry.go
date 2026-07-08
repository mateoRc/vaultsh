package shell

import (
	"sort"

	"github.com/mateom/vaultsh/internal/command"
)

type Registry struct {
	commands map[string]command.Command
}

func NewRegistry() *Registry {
	return &Registry{
		commands: make(map[string]command.Command),
	}
}

func (r *Registry) Register(command command.Command) {
	r.commands[command.Name()] = command
}

func (r *Registry) Find(name string) (command.Command, bool) {
	command, found := r.commands[name]
	return command, found
}

func (r *Registry) Commands() []command.Command {
	commands := make([]command.Command, 0, len(r.commands))
	for _, command := range r.commands {
		commands = append(commands, command)
	}

	sort.Slice(commands, func(i, j int) bool {
		return commands[i].Name() < commands[j].Name()
	})

	return commands
}
