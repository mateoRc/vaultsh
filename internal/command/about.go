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
		Output: "Vaultsh is the read-only virtual shell and Explorer for Backend Lab.\n\n" +
			"Start with /projects/overview.md for the integrated system, " +
			"then /docs/architecture/overview.md for its boundaries. " +
			"Professional background is under /cv.",
		ExitCode: ExitSuccess,
	}
}
