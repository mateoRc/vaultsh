package shell

import (
	"fmt"

	"github.com/mateom/vaultsh/internal/command"
)

type Engine struct {
	commands *command.Registry
}

func New() *Engine {
	commands := command.NewRegistry()
	commands.Register(command.About{})
	commands.Register(command.Clear{})
	commands.Register(command.NewHelp(commands))

	return &Engine{commands: commands}
}

func (e *Engine) Execute(line string) command.Result {
	run, found := e.commands.Find(line)
	if !found {
		return command.Result{
			Output:   fmt.Sprintf("command not found: %s", line),
			ExitCode: 127,
		}
	}

	return run.Execute()
}
