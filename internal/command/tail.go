package command

import "github.com/mateom/vaultsh/internal/filesystem"

const defaultTailLines = 10

type Tail struct {
	workingDirectory *filesystem.WorkingDirectory
}

func NewTail(workingDirectory *filesystem.WorkingDirectory) Tail {
	return Tail{workingDirectory: workingDirectory}
}

func (Tail) Name() string {
	return "tail"
}

func (Tail) Description() string {
	return "Print the last lines"
}

func (Tail) Usage() string {
	return "tail [-n count] [file]"
}

func (t Tail) Execute(args []string, input Input) Result {
	count, path, result := parseLineCount("tail", defaultTailLines, args)
	if result != nil {
		return *result
	}

	data, result := readInput("tail", path, input, t.workingDirectory)
	if result != nil {
		return *result
	}

	currentLines := lines(data)
	if count < len(currentLines) {
		currentLines = currentLines[len(currentLines)-count:]
	}
	return Result{Output: joinLines(currentLines), ExitCode: ExitSuccess}
}
