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
		Output: "I wanted to build a portfolio, so naturally I built a shell first.\n\n" +
			"Vault: a virtual shell wired to search, telemetry, and an Oracle.\n\n" +
			"Browse projects: ls /docs/projects\n" +
			"Ask the Oracle: invoke what is Vault?\n" +
			"Telemetry dashboard: dashboard\n" +
			"All commands: help\n\n" +
			"Switch to Explorer above to browse the documents.",
		ExitCode: ExitSuccess,
	}
}
