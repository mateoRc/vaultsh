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
		Output: "I built my portfolio as a shell.\n\n" +
			"Welcome to Vaultsh.\n\n" +
			"[Contact](mailto:mahmutovic.mateo@gmail.com) · " +
			"[GitHub](https://github.com/mateoRc/vaultsh) · " +
			"[LinkedIn](https://www.linkedin.com/in/mateo-mahmutovi%C4%87-a9837232b/)\n\n" +
			"Explore Mateo's CV & project docs:\n" +
			"  about (project overview)\n" +
			"  whoami (about Mateo)\n" +
			"  tree /cv/experience (browse experience)\n" +
			"  cat /cv/skills.md (view skills)\n" +
			"  search \"Languages\" (search the portfolio)\n" +
			"  contact (email, GitHub, and LinkedIn)\n" +
			"  dashboard (live service activity)\n\n" +
			"Type help or choose a suggestion below.",
		ExitCode: ExitSuccess,
	}
}
