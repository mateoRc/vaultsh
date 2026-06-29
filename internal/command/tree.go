package command

import (
	"fmt"
	"strings"

	"github.com/mateom/vaultsh/internal/filesystem"
)

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
	return "tree [path]"
}

func (t Tree) Execute(args []string, _ Input) Result {
	if len(args) > 1 {
		return Result{
			Output:   "usage: tree [path]",
			ExitCode: ExitUsage,
		}
	}

	target := "."
	if len(args) == 1 {
		target = args[0]
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
	writeChildren(&output, node, "")

	return Result{
		Output:   output.String(),
		ExitCode: ExitSuccess,
	}
}

func writeChildren(output *strings.Builder, node filesystem.Node, prefix string) {
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
		writeChildren(output, child, childPrefix)
	}
}
