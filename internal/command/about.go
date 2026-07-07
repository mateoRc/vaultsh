package command

type About struct{}

func (About) Name() string {
	return "about"
}

func (About) Description() string {
	return "Describe Vaultsh"
}

func (About) Execute([]string, Input) Result {
	return Result{
		Output: "Vaultsh is a read-only virtual shell for exploring Mateo's " +
			"CV, projects, and live backend services.\n\n" +
			"The portfolio is the system: a Go shell backed by " +
			"Atlas search, Forge telemetry, and Lab deployment docs.",
		ExitCode: ExitSuccess,
	}
}
