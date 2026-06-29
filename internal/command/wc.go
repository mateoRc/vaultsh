package command

import (
	"fmt"
	"strings"

	"github.com/mateom/vaultsh/internal/filesystem"
)

type Wc struct {
	workingDirectory *filesystem.WorkingDirectory
}

type wcOptions struct {
	lines bool
	words bool
	bytes bool
	path  string
}

func NewWc(workingDirectory *filesystem.WorkingDirectory) Wc {
	return Wc{workingDirectory: workingDirectory}
}

func (Wc) Name() string {
	return "wc"
}

func (Wc) Description() string {
	return "Count lines, words and bytes"
}

func (Wc) Usage() string {
	return "wc [-lwc] [file]"
}

func (w Wc) Execute(args []string, input Input) Result {
	options, result := parseWcOptions(args)
	if result != nil {
		return *result
	}

	data, result := readInput("wc", options.path, input, w.workingDirectory)
	if result != nil {
		return *result
	}

	var counts []string
	if options.lines {
		counts = append(counts, fmt.Sprint(len(lines(data))))
	}
	if options.words {
		counts = append(counts, fmt.Sprint(len(strings.Fields(data))))
	}
	if options.bytes {
		counts = append(counts, fmt.Sprint(len(data)))
	}

	return Result{
		Output:   strings.Join(counts, " "),
		ExitCode: ExitSuccess,
	}
}

func parseWcOptions(args []string) (wcOptions, *Result) {
	var options wcOptions
	selected := false

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") && arg != "-" {
			for _, option := range strings.TrimPrefix(arg, "-") {
				switch option {
				case 'l':
					options.lines = true
				case 'w':
					options.words = true
				case 'c':
					options.bytes = true
				default:
					return options, &Result{
						Output:   fmt.Sprintf("wc: unknown option -- %c", option),
						ExitCode: ExitUsage,
					}
				}
			}
			selected = true
			continue
		}
		if options.path != "" {
			return options, &Result{
				Output:   "usage: wc [-lwc] [file]",
				ExitCode: ExitUsage,
			}
		}
		options.path = arg
	}

	if !selected {
		options.lines = true
		options.words = true
		options.bytes = true
	}
	return options, nil
}
