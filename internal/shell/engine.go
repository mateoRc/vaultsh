package shell

import (
	"fmt"
	"strings"

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
	commands.Register(command.NewCat(workingDirectory))
	commands.Register(command.NewCd(workingDirectory))
	commands.Register(command.Clear{})
	commands.Register(command.NewHelp(commands))
	commands.Register(command.NewLs(workingDirectory))
	commands.Register(command.NewPwd(workingDirectory))

	return &Engine{commands: commands}
}

func (e *Engine) Execute(line string) command.Result {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return command.Result{}
	}

	run, found := e.commands.Find(fields[0])
	if !found {
		return command.Result{
			Output:   fmt.Sprintf("command not found: %s", fields[0]),
			ExitCode: 127,
		}
	}

	return run.Execute(fields[1:])
}
