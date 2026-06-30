package command

type Hire struct{}

func (Hire) Name() string {
	return "hire"
}

func (Hire) Description() string {
	return "Hire Mateo"
}

func (Hire) Hidden() bool {
	return true
}

func (Hire) Execute(args []string, _ Input) Result {
	if len(args) != 1 || args[0] != "mateo" {
		return Result{
			Output:   "usage: hire mateo",
			ExitCode: ExitUsage,
		}
	}

	return Result{
		Output:   "hire: permission denied\nhint: try sudo hire mateo -s <yearly>",
		ExitCode: ExitFailure,
	}
}
