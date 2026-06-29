package shell

import (
	"fmt"

	"github.com/mateom/vaultsh/internal/command"
)

type Engine struct {
	commands map[string]command.Func
}

func New() *Engine {
	return &Engine{
		commands: map[string]command.Func{
			"about": command.About,
			"help":  command.Help,
		},
	}
}

func (e *Engine) Execute(line string) command.Result {
	run, found := e.commands[line]
	if !found {
		return command.Result{
			Output:   fmt.Sprintf("command not found: %s", line),
			ExitCode: 127,
		}
	}

	return run()
}
