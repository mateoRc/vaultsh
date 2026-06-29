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
	if len(tokens) == 0 {
		return command.Result{}
	}
	e.context.History().Add(line)

	run, found := e.commands.Find(tokens[0].Value)
	if !found {
		return command.Result{
			Output:   fmt.Sprintf("command not found: %s", tokens[0].Value),
			ExitCode: command.ExitNotFound,
		}
	}

	args := make([]string, 0, len(tokens)-1)
	for _, token := range tokens[1:] {
		args = append(args, token.Value)
	}

	return run.Execute(args)
}
