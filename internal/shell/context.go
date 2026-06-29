package shell

import (
	"github.com/mateom/vaultsh/internal/filesystem"
	"github.com/mateom/vaultsh/internal/history"
)

const sessionHistoryLimit = 100

type ExecutionContext struct {
	workingDirectory *filesystem.WorkingDirectory
	history          *history.Store
}

func NewExecutionContext(root *filesystem.Directory) *ExecutionContext {
	return &ExecutionContext{
		workingDirectory: filesystem.NewWorkingDirectory(root),
		history:          history.New(sessionHistoryLimit),
	}
}

func (c *ExecutionContext) WorkingDirectory() *filesystem.WorkingDirectory {
	return c.workingDirectory
}

func (c *ExecutionContext) History() *history.Store {
	return c.history
}
