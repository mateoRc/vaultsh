package command

type Clear struct{}

func (Clear) Name() string {
	return "clear"
}

func (Clear) Description() string {
	return "Clear the terminal"
}

func (Clear) Usage() string {
	return "clear"
}

func (Clear) Execute(args []string, _ Input) Result {
	if len(args) != 0 {
		return Result{Output: "usage: clear", ExitCode: ExitUsage}
	}

	return Result{
		ExitCode: ExitSuccess,
		Action:   ActionClear,
	}
}
