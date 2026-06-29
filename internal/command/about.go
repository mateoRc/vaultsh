package command

type About struct{}

func (About) Name() string {
	return "about"
}

func (About) Description() string {
	return "Describe Vaultsh"
}

func (About) Execute() Result {
	return Result{
		Output:   "Vaultsh is a read-only virtual shell engine.",
		ExitCode: 0,
	}
}
