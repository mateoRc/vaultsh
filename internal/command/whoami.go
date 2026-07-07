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

func (Whoami) Usage() string {
	return "whoami"
}

func (Whoami) Execute(args []string, _ Input) Result {
	if len(args) != 0 {
		return Result{Output: "usage: whoami", ExitCode: ExitUsage}
	}

	return Result{
		Output: "Mateo Mahmutović\n" +
			"Senior Backend Engineer\n" +
			"Currently building distributed backend systems.\n\n" +
			"[Email](mailto:mahmutovic.mateo@gmail.com)\n" +
			"[GitHub](https://github.com/mateoRc)\n" +
			"[LinkedIn](https://www.linkedin.com/in/mateo-mahmutovi%C4%87-a9837232b/)",
		ExitCode: ExitSuccess,
	}
}
