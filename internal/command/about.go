package command

type About struct{}

func (About) Name() string {
	return "about"
}

func (About) Description() string {
	return "Describe Vaultsh"
}

func (About) Usage() string {
	return "about"
}

func (About) Execute(args []string, _ Input) Result {
	if len(args) != 0 {
		return Result{Output: "usage: about", ExitCode: ExitUsage}
	}

	return Result{
		Output: "Vaultsh is a read-only virtual shell for exploring Mateo's " +
			"CV, projects, and live backend services.\n\n" +
			"The portfolio is the system: a Go shell backed by " +
			"Atlas search, Forge telemetry, and Lab deployment docs.",
		ExitCode: ExitSuccess,
	}
}
