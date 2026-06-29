package command

import (
	"fmt"
	"strings"

	"github.com/mateom/vaultsh/internal/filesystem"
)

type Ls struct {
	workingDirectory *filesystem.WorkingDirectory
}

func NewLs(workingDirectory *filesystem.WorkingDirectory) Ls {
	return Ls{workingDirectory: workingDirectory}
}

func (Ls) Name() string {
	return "ls"
}

func (Ls) Description() string {
	return "List directory contents"
}

func (Ls) Usage() string {
	return "ls [path]"
}

func (l Ls) Execute(args []string) Result {
	if len(args) > 1 {
		return Result{
			Output:   "usage: ls [path]",
			ExitCode: ExitUsage,
		}
	}

	target := "."
	if len(args) == 1 {
		target = args[0]
	}

	node, _, err := l.workingDirectory.Resolve(target)
	if err != nil {
		return Result{
			Output:   fmt.Sprintf("ls: %s: %v", target, err),
			ExitCode: ExitFailure,
		}
	}

	directory, ok := node.(*filesystem.Directory)
	if !ok {
		return Result{
			Output:   node.Name(),
			ExitCode: ExitSuccess,
		}
	}

	children := directory.Children()
	names := make([]string, 0, len(children))
	for _, child := range children {
		name := child.Name()
		if child.Kind() == filesystem.KindDirectory {
			name += "/"
		}
		names = append(names, name)
	}

	return Result{
		Output:   strings.Join(names, "\n"),
		ExitCode: ExitSuccess,
	}
}
