package command

import (
	"fmt"
	"strings"

	"github.com/mateom/vaultsh/internal/filesystem"
)

func readInput(
	name string,
	path string,
	input Input,
	workingDirectory *filesystem.WorkingDirectory,
) (string, *Result) {
	if path == "" {
		if input.Present {
			return input.Data, nil
		}
		return "", &Result{
			Output:   fmt.Sprintf("%s: missing input", name),
			ExitCode: ExitUsage,
		}
	}

	node, _, err := workingDirectory.Resolve(path)
	if err != nil {
		return "", &Result{
			Output:   fmt.Sprintf("%s: %s: %v", name, path, err),
			ExitCode: ExitFailure,
		}
	}

	file, ok := node.(*filesystem.File)
	if !ok {
		return "", &Result{
			Output:   fmt.Sprintf("%s: %s: is a directory", name, path),
			ExitCode: ExitFailure,
		}
	}
	return file.Content(), nil
}

func lines(data string) []string {
	if data == "" {
		return nil
	}
	return strings.Split(strings.TrimSuffix(data, "\n"), "\n")
}

func joinLines(values []string) string {
	return strings.Join(values, "\n")
}
