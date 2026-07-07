package command

type Contact struct{}

func (Contact) Name() string {
	return "contact"
}

func (Contact) Description() string {
	return "Show contact links"
}

func (Contact) Usage() string {
	return "contact"
}

func (Contact) Help() string {
	return "In the browser terminal, labels render as clickable links."
}

func (Contact) Execute(args []string, _ Input) Result {
	if len(args) != 0 {
		return Result{Output: "usage: contact", ExitCode: ExitUsage}
	}

	return Result{
		Output: "Contact Mateo:\n" +
			"[Email](mailto:mahmutovic.mateo@gmail.com)\n" +
			"[GitHub](https://github.com/mateoRc)\n" +
			"[LinkedIn](https://www.linkedin.com/in/mateo-mahmutovi%C4%87-a9837232b/)",
		ExitCode: ExitSuccess,
	}
}
