package command

type Action string

const (
	ActionNone  Action = ""
	ActionClear Action = "clear"
)

type Result struct {
	Output   string
	ExitCode int
	Action   Action
}

type Command interface {
	Name() string
	Description() string
	Execute() Result
}
