package command

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mateom/vaultsh/internal/filesystem"
)

type Grep struct {
	workingDirectory *filesystem.WorkingDirectory
}

type grepOptions struct {
	ignoreCase  bool
	lineNumbers bool
	pattern     string
	path        string
}

func NewGrep(workingDirectory *filesystem.WorkingDirectory) Grep {
	return Grep{workingDirectory: workingDirectory}
}

func (Grep) Name() string {
	return "grep"
}

func (Grep) Description() string {
	return "Filter lines by a regular expression"
}

func (Grep) Usage() string {
	return "grep [-in] <pattern> [file]"
}

func (Grep) Help() string {
	return `Example: grep -in "backend" /cv/about.md`
}

func (g Grep) Execute(args []string, input Input) Result {
	options, result := parseGrepOptions(args)
	if result != nil {
		return *result
	}

	data := input.Data
	if options.path != "" {
		node, _, err := g.workingDirectory.Resolve(options.path)
		if err != nil {
			return Result{
				Output:   fmt.Sprintf("grep: %s: %v", options.path, err),
				ExitCode: ExitFailure,
			}
		}

		file, ok := node.(*filesystem.File)
		if !ok {
			return Result{
				Output:   fmt.Sprintf("grep: %s: is a directory", options.path),
				ExitCode: ExitFailure,
			}
		}
		data = file.Content()
	} else if !input.Present {
		return Result{
			Output:   "usage: grep [-in] <pattern> [file]",
			ExitCode: ExitUsage,
		}
	}

	pattern := options.pattern
	if options.ignoreCase {
		pattern = "(?i)" + pattern
	}
	expression, err := regexp.Compile(pattern)
	if err != nil {
		return Result{
			Output:   fmt.Sprintf("grep: invalid pattern: %v", err),
			ExitCode: ExitUsage,
		}
	}

	var matches []string
	for index, line := range lines(data) {
		if !expression.MatchString(line) {
			continue
		}
		if options.lineNumbers {
			line = fmt.Sprintf("%d:%s", index+1, line)
		}
		matches = append(matches, line)
	}

	exitCode := ExitSuccess
	if len(matches) == 0 {
		exitCode = ExitFailure
	}
	return Result{
		Output:   strings.Join(matches, "\n"),
		ExitCode: exitCode,
	}
}

func parseGrepOptions(args []string) (grepOptions, *Result) {
	var options grepOptions
	optionsEnded := false
	var positional []string

	for _, arg := range args {
		if !optionsEnded && arg == "--" {
			optionsEnded = true
			continue
		}
		if !optionsEnded && strings.HasPrefix(arg, "-") && arg != "-" {
			for _, option := range strings.TrimPrefix(arg, "-") {
				switch option {
				case 'i':
					options.ignoreCase = true
				case 'n':
					options.lineNumbers = true
				default:
					return options, &Result{
						Output:   fmt.Sprintf("grep: unknown option -- %c", option),
						ExitCode: ExitUsage,
					}
				}
			}
			continue
		}
		positional = append(positional, arg)
	}

	if len(positional) < 1 || len(positional) > 2 {
		return options, &Result{
			Output:   "usage: grep [-in] <pattern> [file]",
			ExitCode: ExitUsage,
		}
	}

	options.pattern = positional[0]
	if len(positional) == 2 {
		options.path = positional[1]
	}
	return options, nil
}
