package command

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mateom/vaultsh/internal/filesystem"
)

const unlimitedTreeDepth = -1

type Tree struct {
	workingDirectory *filesystem.WorkingDirectory
}

func NewTree(workingDirectory *filesystem.WorkingDirectory) Tree {
	return Tree{workingDirectory: workingDirectory}
}

func (Tree) Name() string {
	return "tree"
}

func (Tree) Description() string {
	return "Print a directory tree"
}

func (Tree) Usage() string {
	return "tree [-L depth] [path]"
}

func (t Tree) Execute(args []string, _ Input) Result {
	depth, target, result := parseTreeOptions(args)
	if result != nil {
		return *result
	}

	node, resolvedPath, err := t.workingDirectory.Resolve(target)
	if err != nil {
		return Result{
			Output:   fmt.Sprintf("tree: %s: %v", target, err),
			ExitCode: ExitFailure,
		}
	}

	var output strings.Builder
	label := node.Name()
	if target == "." {
		label = "."
	} else if resolvedPath == "/" {
		label = "/"
	}
	output.WriteString(label)
	writeChildren(&output, node, "", depth, 0)

	return Result{
		Output:   output.String(),
		ExitCode: ExitSuccess,
	}
}

func parseTreeOptions(args []string) (int, string, *Result) {
	depth := unlimitedTreeDepth
	target := "."
	pathSet := false

	for index := 0; index < len(args); index++ {
		if args[index] == "-L" {
			if index+1 >= len(args) {
				return 0, "", treeUsage()
			}
			value, err := strconv.Atoi(args[index+1])
			if err != nil || value < 1 {
				return 0, "", &Result{
					Output:   fmt.Sprintf("tree: invalid depth: %s", args[index+1]),
					ExitCode: ExitUsage,
				}
			}
			depth = value
			index++
			continue
		}
		if strings.HasPrefix(args[index], "-") && args[index] != "-" {
			return 0, "", &Result{
				Output:   fmt.Sprintf("tree: unknown option: %s", args[index]),
				ExitCode: ExitUsage,
			}
		}
		if pathSet {
			return 0, "", treeUsage()
		}
		target = args[index]
		pathSet = true
	}

	return depth, target, nil
}

func treeUsage() *Result {
	return &Result{
		Output:   "usage: tree [-L depth] [path]",
		ExitCode: ExitUsage,
	}
}

func writeChildren(
	output *strings.Builder,
	node filesystem.Node,
	prefix string,
	maxDepth int,
	depth int,
) {
	if maxDepth != unlimitedTreeDepth && depth >= maxDepth {
		return
	}

	directory, ok := node.(*filesystem.Directory)
	if !ok {
		return
	}

	children := directory.Children()
	for index, child := range children {
		last := index == len(children)-1
		connector := "├── "
		childPrefix := prefix + "│   "
		if last {
			connector = "└── "
			childPrefix = prefix + "    "
		}

		fmt.Fprintf(output, "\n%s%s%s", prefix, connector, child.Name())
		writeChildren(output, child, childPrefix, maxDepth, depth+1)
	}
}
