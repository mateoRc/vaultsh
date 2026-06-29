package shell

import "fmt"

type Result struct {
	Output   string
	ExitCode int
}

type commandFunc func() Result

type Engine struct {
	commands map[string]commandFunc
}

func New() *Engine {
	return &Engine{
		commands: map[string]commandFunc{
			"about": about,
			"help": help,
		},
	}
}

func (e *Engine) Execute(line string) Result {
	command, found := e.commands[line]
	if !found {
		return Result{
			Output:   fmt.Sprintf("command not found: %s", line),
			ExitCode: 127,
		}
	}

	return command()
}
