package shell

import (
	"fmt"
	"strings"

	"github.com/mateom/vaultsh/internal/command"
	"github.com/mateom/vaultsh/internal/filesystem"
)

type Engine struct {
	commands *command.Registry
	context  *ExecutionContext
}

func New() *Engine {
	return NewWithRoot(filesystem.NewDirectory(""))
}

func NewWithRoot(root *filesystem.Directory) *Engine {
	return NewWithContext(NewExecutionContext(root))
}

func NewWithContext(context *ExecutionContext) *Engine {
	commands := command.NewRegistry()
	workingDirectory := context.WorkingDirectory()

	commands.Register(command.About{})
	commands.Register(command.NewCat(workingDirectory))
	commands.Register(command.NewCd(workingDirectory))
	commands.Register(command.Clear{})
	commands.Register(command.NewHelp(commands))
	commands.Register(command.NewHistory(context.History()))
	commands.Register(command.NewLs(workingDirectory))
	commands.Register(command.NewPwd(workingDirectory))
	commands.Register(command.NewTree(workingDirectory))

	return &Engine{
		commands: commands,
		context:  context,
	}
}

func (e *Engine) Execute(line string) command.Result {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return command.Result{}
	}
	e.context.History().Add(line)

	run, found := e.commands.Find(fields[0])
	if !found {
		return command.Result{
			Output:   fmt.Sprintf("command not found: %s", fields[0]),
			ExitCode: 127,
		}
	}

	return run.Execute(fields[1:])
}
