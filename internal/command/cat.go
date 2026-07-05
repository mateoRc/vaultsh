package command

import (
	"fmt"
	"strings"

	"github.com/mateom/vaultsh/internal/filesystem"
)

type Cat struct {
	workingDirectory *filesystem.WorkingDirectory
}

func NewCat(workingDirectory *filesystem.WorkingDirectory) Cat {
	return Cat{workingDirectory: workingDirectory}
}

func (Cat) Name() string {
	return "cat"
}

func (Cat) Description() string {
	return "Print file contents"
}

func (Cat) Usage() string {
	return "cat [-n] [file]"
}

func (Cat) Help() string {
	return "Example: cat /cv/about.md"
}

func (c Cat) Execute(args []string, input Input) Result {
	numberLines, path, result := parseCatOptions(args)
	if result != nil {
		return *result
	}
	if path == "" && !input.Present {
		return Result{
			Output:   "usage: cat [-n] [file]",
			ExitCode: ExitUsage,
		}
	}

	data, result := readInput("cat", path, input, c.workingDirectory)
	if result != nil {
		return *result
	}

	if numberLines {
		numbered := make([]string, 0, len(lines(data)))
		for index, line := range lines(data) {
			numbered = append(numbered, fmt.Sprintf("%6d\t%s", index+1, line))
		}
		data = joinLines(numbered)
	}

	return Result{
		Output:   data,
		ExitCode: ExitSuccess,
	}
}

func parseCatOptions(args []string) (bool, string, *Result) {
	numberLines := false
	optionsEnded := false
	var path string

	for _, arg := range args {
		if !optionsEnded && arg == "--" {
			optionsEnded = true
			continue
		}
		if !optionsEnded && arg == "-n" {
			numberLines = true
			continue
		}
		if !optionsEnded && strings.HasPrefix(arg, "-") && arg != "-" {
			return false, "", &Result{
				Output:   fmt.Sprintf("cat: unknown option: %s", arg),
				ExitCode: ExitUsage,
			}
		}
		if path != "" {
			return false, "", &Result{
				Output:   "usage: cat [-n] [file]",
				ExitCode: ExitUsage,
			}
		}
		path = arg
	}

	return numberLines, path, nil
}
