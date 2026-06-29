package shell

import "github.com/mateom/vaultsh/internal/filesystem"

type ExecutionContext struct {
	workingDirectory *filesystem.WorkingDirectory
}

func NewExecutionContext(root *filesystem.Directory) *ExecutionContext {
	return &ExecutionContext{
		workingDirectory: filesystem.NewWorkingDirectory(root),
	}
}

func (c *ExecutionContext) WorkingDirectory() *filesystem.WorkingDirectory {
	return c.workingDirectory
}
