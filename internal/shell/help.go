package shell

func help() Result {
	return Result{
		Output:   "Available commands:\n  about\n  help",
		ExitCode: 0,
	}
}
