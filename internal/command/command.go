package command

type Action string

const (
	ActionNone  Action = ""
	ActionClear Action = "clear"
)

const (
	ExitSuccess     = 0
	ExitFailure     = 1
	ExitUsage       = 2
	ExitUnsupported = 126
	ExitNotFound    = 127
)

type Result struct {
	Output   string
	ExitCode int
	Action   Action
	Verbose  string
}

type Input struct {
	Data    string
	Present bool
}

type Command interface {
	Name() string
	Description() string
	Execute(args []string, input Input) Result
}

func IsHidden(current Command) bool {
	hidden, ok := current.(interface{ Hidden() bool })
	return ok && hidden.Hidden()
}
