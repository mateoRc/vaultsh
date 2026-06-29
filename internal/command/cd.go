package command

import (
	"fmt"

	"github.com/mateom/vaultsh/internal/filesystem"
)

type Cd struct {
	workingDirectory *filesystem.WorkingDirectory
}

func NewCd(workingDirectory *filesystem.WorkingDirectory) Cd {
	return Cd{workingDirectory: workingDirectory}
}

func (Cd) Name() string {
	return "cd"
}

func (Cd) Description() string {
	return "Change the current directory"
}

func (Cd) Usage() string {
	return "cd [directory]"
}

func (c Cd) Execute(args []string, _ Input) Result {
	if len(args) > 1 {
		return Result{
			Output:   "usage: cd [directory]",
			ExitCode: ExitUsage,
		}
	}

	target := "/"
	if len(args) == 1 {
		target = args[0]
	}

	if err := c.workingDirectory.Change(target); err != nil {
		return Result{
			Output:   fmt.Sprintf("cd: %s: %v", target, err),
			ExitCode: ExitFailure,
		}
	}

	return Result{ExitCode: ExitSuccess}
}
