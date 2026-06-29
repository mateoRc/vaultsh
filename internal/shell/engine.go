package shell

import (
	"fmt"

	"github.com/mateom/vaultsh/internal/command"
	"github.com/mateom/vaultsh/internal/filesystem"
	"github.com/mateom/vaultsh/internal/parser"
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
	tokens, err := parser.Lex(line)
	if err != nil {
		return command.Result{
			Output:   fmt.Sprintf("syntax error: %v", err),
			ExitCode: command.ExitUsage,
		}
	}
	commands, err := parser.Parse(tokens)
	if err != nil {
		return command.Result{
			Output:   fmt.Sprintf("syntax error: %v", err),
			ExitCode: command.ExitUsage,
		}
	}
	if len(commands) == 0 {
		return command.Result{}
	}
	if len(commands) > 1 {
		return command.Result{
			Output:   "pipelines are not supported yet",
			ExitCode: command.ExitUsage,
		}
	}
	e.context.History().Add(line)

	current := commands[0]
	run, found := e.commands.Find(current[0])
	if !found {
		return command.Result{
			Output:   fmt.Sprintf("command not found: %s", current[0]),
			ExitCode: command.ExitNotFound,
		}
	}

	return run.Execute(current[1:])
}
