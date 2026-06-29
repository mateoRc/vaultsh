package command

import (
	"fmt"
	"strconv"

	"github.com/mateom/vaultsh/internal/filesystem"
)

const defaultHeadLines = 10

type Head struct {
	workingDirectory *filesystem.WorkingDirectory
}

func NewHead(workingDirectory *filesystem.WorkingDirectory) Head {
	return Head{workingDirectory: workingDirectory}
}

func (Head) Name() string {
	return "head"
}

func (Head) Description() string {
	return "Print the first lines"
}

func (Head) Usage() string {
	return "head [-n count] [file]"
}

func (h Head) Execute(args []string, input Input) Result {
	count, path, result := parseLineCount("head", defaultHeadLines, args)
	if result != nil {
		return *result
	}

	data, result := readInput("head", path, input, h.workingDirectory)
	if result != nil {
		return *result
	}

	currentLines := lines(data)
	if count < len(currentLines) {
		currentLines = currentLines[:count]
	}
	return Result{Output: joinLines(currentLines), ExitCode: ExitSuccess}
}

func parseLineCount(name string, defaultCount int, args []string) (int, string, *Result) {
	count := defaultCount
	var path string

	for index := 0; index < len(args); index++ {
		if args[index] == "-n" {
			if index+1 >= len(args) {
				return 0, "", lineCountUsage(name)
			}
			value, err := strconv.Atoi(args[index+1])
			if err != nil || value < 0 {
				return 0, "", &Result{
					Output:   fmt.Sprintf("%s: invalid line count: %s", name, args[index+1]),
					ExitCode: ExitUsage,
				}
			}
			count = value
			index++
			continue
		}
		if path != "" {
			return 0, "", lineCountUsage(name)
		}
		path = args[index]
	}

	return count, path, nil
}

func lineCountUsage(name string) *Result {
	return &Result{
		Output:   fmt.Sprintf("usage: %s [-n count] [file]", name),
		ExitCode: ExitUsage,
	}
}
