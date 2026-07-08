package shell

import (
	"fmt"
	"strings"

	"github.com/mateom/vaultsh/internal/command"
	"github.com/mateom/vaultsh/internal/filesystem"
	"github.com/mateom/vaultsh/internal/parser"
)

type Engine struct {
	commands *Registry
	context  *ExecutionContext
}

type Dependencies struct {
	Search      command.SearchService
	Metrics     command.MetricsService
	Deployments command.DeploymentService
	System      command.SystemService
	Assessment  command.AssessmentService
	Events      EventRecorder
}

type EventRecorder interface {
	Record(service, event, name string, durationMS int64, exitCode int) error
}

func New() *Engine {
	return NewWithRoot(filesystem.NewDirectory(""))
}

func NewWithRoot(root *filesystem.Directory) *Engine {
	return NewWithContext(NewExecutionContext(root))
}

func NewWithContext(context *ExecutionContext) *Engine {
	return NewWithContextAndDependencies(context, Dependencies{})
}

func NewWithContextAndDependencies(
	context *ExecutionContext,
	dependencies Dependencies,
) *Engine {
	commands := NewRegistry()
	workingDirectory := context.WorkingDirectory()
	negotiation := context.Negotiation()

	commands.Register(NewAcceptOffer("Y", negotiation))
	commands.Register(command.About{})
	commands.Register(command.NewCat(workingDirectory))
	commands.Register(command.NewCd(workingDirectory))
	commands.Register(command.Clear{})
	commands.Register(command.Contact{})
	commands.Register(command.NewGrep(workingDirectory))
	commands.Register(command.NewHead(workingDirectory))
	commands.Register(command.NewHelp(commands))
	commands.Register(command.NewHistory(context.History()))
	commands.Register(command.Hire{})
	commands.Register(command.NewLs(workingDirectory))
	commands.Register(command.NewPwd(workingDirectory))
	commands.Register(command.NewSort(workingDirectory))
	commands.Register(command.NewSudo(negotiation))
	commands.Register(command.NewTail(workingDirectory))
	commands.Register(command.NewTree(workingDirectory))
	commands.Register(command.NewWc(workingDirectory))
	commands.Register(command.Welcome{})
	commands.Register(command.Whoami{})
	commands.Register(NewAcceptOffer("y", negotiation))
	if dependencies.Search != nil {
		commands.Register(command.NewSearch(dependencies.Search))
	}
	if dependencies.Metrics != nil {
		commands.Register(command.NewMetrics(dependencies.Metrics))
		commands.Register(command.NewDashboard(
			dependencies.Metrics,
			dependencies.Deployments,
			dependencies.System,
			dependencies.Assessment,
		))
	}
	if dependencies.Deployments != nil {
		commands.Register(command.NewDeployments(dependencies.Deployments))
	}

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
	syntaxTree, err := parser.Parse(tokens)
	if err != nil {
		return command.Result{
			Output:   fmt.Sprintf("syntax error: %v", err),
			ExitCode: command.ExitUsage,
		}
	}
	if len(syntaxTree.Pipeline) == 0 {
		return command.Result{}
	}
	verbose := takeVerboseFlag(&syntaxTree)
	e.context.History().Add(line)

	var result command.Result
	input := command.Input{}
	completed := 0
	for _, current := range syntaxTree.Pipeline {
		run, found := e.commands.Find(current.Name)
		if !found {
			result = command.Result{
				Output:   fmt.Sprintf("command not found: %s", current.Name),
				ExitCode: command.ExitNotFound,
			}
			return withVerbose(result, syntaxTree, completed, verbose)
		}

		result = run.Execute(current.Args, input)
		if result.ExitCode != command.ExitSuccess {
			return withVerbose(result, syntaxTree, completed, verbose)
		}
		completed++
		input = command.Input{
			Data:    result.Output,
			Present: true,
		}
	}

	return withVerbose(result, syntaxTree, completed, verbose)
}

func takeVerboseFlag(tree *parser.SyntaxTree) bool {
	last := &tree.Pipeline[len(tree.Pipeline)-1]
	if len(last.Args) == 0 || last.Args[len(last.Args)-1] != "--verbose" {
		return false
	}
	last.Args = last.Args[:len(last.Args)-1]
	return true
}

func withVerbose(
	result command.Result,
	tree parser.SyntaxTree,
	completed int,
	verbose bool,
) command.Result {
	if !verbose {
		return result
	}

	names := make([]string, len(tree.Pipeline))
	for index, current := range tree.Pipeline {
		names[index] = current.Name
	}
	result.Verbose = fmt.Sprintf(
		"pipeline=%s; stages=%d; completed=%d",
		strings.Join(names, ","),
		len(names),
		completed,
	)
	return result
}
