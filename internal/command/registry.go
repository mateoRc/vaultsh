package command

import "sort"

type Registry struct {
	commands map[string]Command
}

func NewRegistry() *Registry {
	return &Registry{
		commands: make(map[string]Command),
	}
}

func (r *Registry) Register(command Command) {
	r.commands[command.Name()] = command
}

func (r *Registry) Find(name string) (Command, bool) {
	command, found := r.commands[name]
	return command, found
}

func (r *Registry) Commands() []Command {
	commands := make([]Command, 0, len(r.commands))
	for _, command := range r.commands {
		commands = append(commands, command)
	}

	sort.Slice(commands, func(i, j int) bool {
		return commands[i].Name() < commands[j].Name()
	})

	return commands
}
