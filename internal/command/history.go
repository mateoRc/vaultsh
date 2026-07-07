package command

import (
	"fmt"
	"strings"

	"github.com/mateom/vaultsh/internal/history"
)

type History struct {
	store *history.Store
}

func NewHistory(store *history.Store) History {
	return History{store: store}
}

func (History) Name() string {
	return "history"
}

func (History) Description() string {
	return "List commands from this session"
}

func (History) Usage() string {
	return "history"
}

func (h History) Execute(args []string, _ Input) Result {
	if len(args) != 0 {
		return Result{Output: "usage: history", ExitCode: ExitUsage}
	}

	var output strings.Builder
	for index, line := range h.store.Entries() {
		if index > 0 {
			output.WriteByte('\n')
		}
		fmt.Fprintf(&output, "%d  %s", index+1, line)
	}

	return Result{
		Output:   output.String(),
		ExitCode: ExitSuccess,
	}
}
