package shell

import "fmt"

type Result struct {
	Output   string
	ExitCode int
}

type Engine struct{}

func (e *Engine) Execute(line string) Result {
	if line != "help" {
		return Result{
			Output:   fmt.Sprintf("command not found: %s", line),
			ExitCode: 127,
		}
	}

	return Result{
		Output:   "Available commands:\n  help",
		ExitCode: 0,
	}
}
