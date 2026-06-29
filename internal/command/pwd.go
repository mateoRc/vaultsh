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

func (p Pwd) Execute() Result {
	return Result{
		Output:   p.workingDirectory.Path(),
		ExitCode: 0,
	}
}
