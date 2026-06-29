package command

import (
	"fmt"

	"github.com/mateom/vaultsh/internal/filesystem"
)

type Cat struct {
	workingDirectory *filesystem.WorkingDirectory
}

func NewCat(workingDirectory *filesystem.WorkingDirectory) Cat {
	return Cat{workingDirectory: workingDirectory}
}

func (Cat) Name() string {
	return "cat"
}

func (Cat) Description() string {
	return "Print file contents"
}

func (Cat) Usage() string {
	return "cat <file>"
}

func (c Cat) Execute(args []string) Result {
	if len(args) != 1 {
		return Result{
			Output:   "usage: cat <file>",
			ExitCode: ExitUsage,
		}
	}

	node, _, err := c.workingDirectory.Resolve(args[0])
	if err != nil {
		return Result{
			Output:   fmt.Sprintf("cat: %s: %v", args[0], err),
			ExitCode: ExitFailure,
		}
	}

	file, ok := node.(*filesystem.File)
	if !ok {
		return Result{
			Output:   fmt.Sprintf("cat: %s: is a directory", args[0]),
			ExitCode: ExitFailure,
		}
	}

	return Result{
		Output:   file.Content(),
		ExitCode: ExitSuccess,
	}
}
