package command

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mateom/vaultsh/internal/filesystem"
)

type Sort struct {
	workingDirectory *filesystem.WorkingDirectory
}

func NewSort(workingDirectory *filesystem.WorkingDirectory) Sort {
	return Sort{workingDirectory: workingDirectory}
}

func (Sort) Name() string {
	return "sort"
}

func (Sort) Description() string {
	return "Sort lines"
}

func (Sort) Usage() string {
	return "sort [-r] [file]"
}

func (s Sort) Execute(args []string, input Input) Result {
	reverse := false
	var path string
	optionsEnded := false
	for _, arg := range args {
		switch {
		case !optionsEnded && arg == "--":
			optionsEnded = true
		case !optionsEnded && arg == "-r":
			reverse = true
		case !optionsEnded && strings.HasPrefix(arg, "-") && arg != "-":
			return Result{
				Output:   fmt.Sprintf("sort: unknown option: %s", arg),
				ExitCode: ExitUsage,
			}
		case path == "":
			path = arg
		default:
			return Result{
				Output:   "usage: sort [-r] [file]",
				ExitCode: ExitUsage,
			}
		}
	}

	data, result := readInput("sort", path, input, s.workingDirectory)
	if result != nil {
		return *result
	}

	currentLines := lines(data)
	sort.Strings(currentLines)
	if reverse {
		for left, right := 0, len(currentLines)-1; left < right; left, right = left+1, right-1 {
			currentLines[left], currentLines[right] = currentLines[right], currentLines[left]
		}
	}
	return Result{Output: joinLines(currentLines), ExitCode: ExitSuccess}
}
