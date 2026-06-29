package command

type Clear struct{}

func (Clear) Name() string {
	return "clear"
}

func (Clear) Description() string {
	return "Clear the terminal"
}

func (Clear) Execute([]string) Result {
	return Result{
		ExitCode: ExitSuccess,
		Action:   ActionClear,
	}
}
