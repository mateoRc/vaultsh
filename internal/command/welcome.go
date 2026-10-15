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
			"Vaultsh: a virtual shell connected to search, telemetry, oracle AI agent.\n\n" +
			"Browse projects: ls /projects\n" +
			"Search: help search\n\n" +
			"Ask oracle anything: invoke hello\n\n" +
			"Telemetry: dashboard"
			"Type help for commands, or switch to Explorer above.",
		ExitCode: ExitSuccess,
	}
}
