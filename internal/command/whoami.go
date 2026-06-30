package command

type Whoami struct{}

func (Whoami) Name() string {
	return "whoami"
}

func (Whoami) Description() string {
	return "Show the current identity"
}

func (Whoami) Hidden() bool {
	return true
}

func (Whoami) Execute([]string, Input) Result {
	return Result{
		Output: "Mateo Mahmutović\n" +
			"Senior Backend Engineer\n" +
			"Currently building distributed backend systems.",
		ExitCode: ExitSuccess,
	}
}
