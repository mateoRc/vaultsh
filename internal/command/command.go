package command

type Result struct {
	Output   string
	ExitCode int
}

type Command interface {
	Name() string
	Description() string
	Execute() Result
}
