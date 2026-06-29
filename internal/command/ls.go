package command

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mateom/vaultsh/internal/filesystem"
)

const (
	directoryMode = "dr-xr-xr-x"
	fileMode      = "-r--r--r--"
)

type Ls struct {
	workingDirectory *filesystem.WorkingDirectory
}

type lsOptions struct {
	all     bool
	long    bool
	path    string
	pathSet bool
}

func NewLs(workingDirectory *filesystem.WorkingDirectory) Ls {
	return Ls{workingDirectory: workingDirectory}
}

func (Ls) Name() string {
	return "ls"
}

func (Ls) Description() string {
	return "List directory contents"
}

func (Ls) Usage() string {
	return "ls [-al] [path]"
}

func (l Ls) Execute(args []string, _ Input) Result {
	options, result := parseLsOptions(args)
	if result != nil {
		return *result
	}

	node, _, err := l.workingDirectory.Resolve(options.path)
	if err != nil {
		return Result{
			Output:   fmt.Sprintf("ls: %s: %v", options.path, err),
			ExitCode: ExitFailure,
		}
	}

	directory, ok := node.(*filesystem.Directory)
	if !ok {
		output := node.Name()
		if options.long {
			output = formatLongEntry(node)
		}
		return Result{
			Output:   output,
			ExitCode: ExitSuccess,
		}
	}

	children := directory.Children()
	names := make([]string, 0, len(children))
	for _, child := range children {
		if !options.all && strings.HasPrefix(child.Name(), ".") {
			continue
		}

		if options.long {
			names = append(names, formatLongEntry(child))
			continue
		}

		name := child.Name()
		if child.Kind() == filesystem.KindDirectory {
			name += "/"
		}
		names = append(names, name)
	}

	return Result{
		Output:   strings.Join(names, "\n"),
		ExitCode: ExitSuccess,
	}
}

func parseLsOptions(args []string) (lsOptions, *Result) {
	options := lsOptions{path: "."}
	optionsEnded := false

	for _, arg := range args {
		if !optionsEnded && arg == "--" {
			optionsEnded = true
			continue
		}
		if !optionsEnded && strings.HasPrefix(arg, "-") && arg != "-" {
			for _, option := range strings.TrimPrefix(arg, "-") {
				switch option {
				case 'a':
					options.all = true
				case 'l':
					options.long = true
				case 't':
					return options, &Result{
						Output:   "ls: option -t requires file timestamps",
						ExitCode: ExitUnsupported,
					}
				default:
					return options, &Result{
						Output:   fmt.Sprintf("ls: unknown option -- %c", option),
						ExitCode: ExitUsage,
					}
				}
			}
			continue
		}
		if options.pathSet {
			return options, &Result{
				Output:   "usage: ls [-al] [path]",
				ExitCode: ExitUsage,
			}
		}
		options.path = arg
		options.pathSet = true
	}

	return options, nil
}

func formatLongEntry(node filesystem.Node) string {
	if directory, ok := node.(*filesystem.Directory); ok {
		return fmt.Sprintf("%s %8s %s/", directoryMode, "-", directory.Name())
	}

	file := node.(*filesystem.File)
	return fmt.Sprintf(
		"%s %8s %s",
		fileMode,
		strconv.Itoa(len(file.Content())),
		file.Name(),
	)
}
