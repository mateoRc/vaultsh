package command

type Welcome struct{}

func (Welcome) Name() string {
	return "welcome"
}

func (Welcome) Description() string {
	return "Show the terminal introduction"
}

func (Welcome) Usage() string {
	return "welcome"
}

func (Welcome) Execute(args []string, _ Input) Result {
	if len(args) != 0 {
		return Result{Output: "usage: welcome", ExitCode: ExitUsage}
	}

	return Result{
		Output: "I built a portfolio, so naturally I built a shell first.\n\n" +
			"Vaultsh opens Backend Lab: a virtual shell connected to search, telemetry, and deployment docs.\n\n" +
			"Explore projects: head -n 11 /projects/overview.md\n" +
			"View architecture: head -n 26 /docs/architecture/overview.md\n" +
			"Live dashboard: dashboard\n\n" +
			"Type help for commands, or switch to Explorer above.",
		ExitCode: ExitSuccess,
	}
}
