package command

import "github.com/mateom/vaultsh/internal/filesystem"

type Pwd struct {
	workingDirectory *filesystem.WorkingDirectory
}

func NewPwd(workingDirectory *filesystem.WorkingDirectory) Pwd {
	return Pwd{workingDirectory: workingDirectory}
}

func (Pwd) Name() string {
	return "pwd"
}

func (Pwd) Description() string {
	return "Print the current directory"
}

func (Pwd) Usage() string {
	return "pwd"
}

func (p Pwd) Execute(args []string, _ Input) Result {
	if len(args) != 0 {
		return Result{Output: "usage: pwd", ExitCode: ExitUsage}
	}

	return Result{
		Output:   p.workingDirectory.Path(),
		ExitCode: ExitSuccess,
	}
}
