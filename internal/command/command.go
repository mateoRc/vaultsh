package command

type Result struct {
	Output   string
	ExitCode int
}

type Func func() Result
