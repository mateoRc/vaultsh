package command

import (
	"fmt"
	"path"
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
	all       bool
	long      bool
	recursive bool
	path      string
	pathSet   bool
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
	return "ls [-alR] [path]"
}

func (l Ls) Execute(args []string, _ Input) Result {
	options, result := parseLsOptions(args)
	if result != nil {
		return *result
	}

	node, resolvedPath, err := l.workingDirectory.Resolve(options.path)
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

	if options.recursive {
		var sections []string
		collectRecursiveListings(directory, resolvedPath, options, &sections)
		return Result{
			Output:   strings.Join(sections, "\n\n"),
			ExitCode: ExitSuccess,
		}
	}

	return Result{
		Output:   strings.Join(formatDirectoryEntries(directory, options), "\n"),
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
				case 'R':
					options.recursive = true
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
				Output:   "usage: ls [-alR] [path]",
				ExitCode: ExitUsage,
			}
		}
		options.path = arg
		options.pathSet = true
	}

	return options, nil
}

func collectRecursiveListings(
	directory *filesystem.Directory,
	directoryPath string,
	options lsOptions,
	sections *[]string,
) {
	entries := visibleChildren(directory, options.all)
	lines := []string{directoryPath + ":"}
	for _, entry := range entries {
		lines = append(lines, formatLsEntry(entry, options.long))
	}
	*sections = append(*sections, strings.Join(lines, "\n"))

	for _, entry := range entries {
		child, ok := entry.(*filesystem.Directory)
		if !ok {
			continue
		}
		collectRecursiveListings(
			child,
			path.Join(directoryPath, child.Name()),
			options,
			sections,
		)
	}
}

func formatDirectoryEntries(directory *filesystem.Directory, options lsOptions) []string {
	children := visibleChildren(directory, options.all)
	entries := make([]string, 0, len(children))
	for _, child := range children {
		entries = append(entries, formatLsEntry(child, options.long))
	}
	return entries
}

func visibleChildren(directory *filesystem.Directory, all bool) []filesystem.Node {
	var children []filesystem.Node
	for _, child := range directory.Children() {
		if !all && strings.HasPrefix(child.Name(), ".") {
			continue
		}
		children = append(children, child)
	}
	return children
}

func formatLsEntry(node filesystem.Node, long bool) string {
	if long {
		return formatLongEntry(node)
	}

	name := node.Name()
	if node.Kind() == filesystem.KindDirectory {
		name += "/"
	}
	return name
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
