package command

import (
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

func (l Ls) Execute() Result {
	children := l.workingDirectory.Directory().Children()
	names := make([]string, 0, len(children))
	for _, child := range children {
		names = append(names, child.Name())
	}

	return Result{
		Output:   strings.Join(names, "\n"),
		ExitCode: 0,
	}
}
