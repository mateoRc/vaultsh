package shell

import (
	"fmt"

	"github.com/mateom/vaultsh/internal/command"
	"github.com/mateom/vaultsh/internal/filesystem"
)

type Engine struct {
	commands *command.Registry
}

func New() *Engine {
	return NewWithRoot(filesystem.NewDirectory(""))
}

func NewWithRoot(root *filesystem.Directory) *Engine {
	commands := command.NewRegistry()
	workingDirectory := filesystem.NewWorkingDirectory(root)

	commands.Register(command.About{})
	commands.Register(command.Clear{})
	commands.Register(command.NewHelp(commands))
	commands.Register(command.NewLs(workingDirectory))
	commands.Register(command.NewPwd(workingDirectory))

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
