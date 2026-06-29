package command

func About() Result {
	return Result{
		Output:   "Vaultsh is a read-only virtual shell engine.",
		ExitCode: 0,
	}
}
