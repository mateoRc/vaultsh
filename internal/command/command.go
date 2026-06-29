package command

type Action string

const (
	ActionNone  Action = ""
	ActionClear Action = "clear"
)

const (
	ExitSuccess  = 0
	ExitFailure  = 1
	ExitUsage    = 2
	ExitNotFound = 127
)

type Result struct {
	Output   string
	ExitCode int
	Action   Action
}

type Command interface {
	Name() string
	Description() string
	Execute(args []string) Result
}
