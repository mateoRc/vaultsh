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

func (h History) Execute([]string) Result {
	var output strings.Builder
	for index, line := range h.store.Entries() {
		if index > 0 {
			output.WriteByte('\n')
		}
		fmt.Fprintf(&output, "%d  %s", index+1, line)
	}

	return Result{
		Output:   output.String(),
		ExitCode: 0,
	}
}
