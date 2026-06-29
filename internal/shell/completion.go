package shell

import (
	"path"
	"sort"
	"strings"

	"github.com/mateom/vaultsh/internal/command"
	"github.com/mateom/vaultsh/internal/filesystem"
)

type Completion struct {
	Start       int
	End         int
	Replacement string
	Candidates  []string
}

func (e *Engine) Complete(line string, cursor int) Completion {
	if cursor < 0 {
		cursor = 0
	}
	if cursor > len(line) {
		cursor = len(line)
	}

	prefix := line[:cursor]
	start := strings.LastIndexAny(prefix, " \t\n|") + 1
	fragment := prefix[start:]
	before := strings.TrimSpace(prefix[:start])

	var candidates []string
	if before == "" || strings.HasSuffix(before, "|") {
		candidates = e.completeCommand(fragment)
	} else {
		candidates = e.completePath(fragment)
	}

	replacement := longestCommonPrefix(candidates)
	if len(candidates) == 1 {
		if strings.HasSuffix(candidates[0], "/") {
			replacement = candidates[0]
		} else {
			replacement = candidates[0] + " "
		}
	}

	return Completion{
		Start:       start,
		End:         cursor,
		Replacement: replacement,
		Candidates:  candidates,
	}
}

func (e *Engine) completeCommand(fragment string) []string {
	var candidates []string
	for _, current := range e.commands.Commands() {
		if command.IsHidden(current) {
			continue
		}
		if strings.HasPrefix(current.Name(), fragment) {
			candidates = append(candidates, current.Name())
		}
	}
	return candidates
}

func (e *Engine) completePath(fragment string) []string {
	directoryPath, namePrefix := path.Split(fragment)
	target := directoryPath
	if target == "" {
		target = "."
	}

	node, _, err := e.context.WorkingDirectory().Resolve(target)
	if err != nil {
		return nil
	}
	directory, ok := node.(*filesystem.Directory)
	if !ok {
		return nil
	}

	var candidates []string
	for _, child := range directory.Children() {
		if !strings.HasPrefix(child.Name(), namePrefix) {
			continue
		}

		name := directoryPath + child.Name()
		if child.Kind() == filesystem.KindDirectory {
			name += "/"
		}
		candidates = append(candidates, name)
	}
	sort.Strings(candidates)
	return candidates
}

func longestCommonPrefix(values []string) string {
	if len(values) == 0 {
		return ""
	}

	prefix := values[0]
	for _, value := range values[1:] {
		for !strings.HasPrefix(value, prefix) {
			prefix = prefix[:len(prefix)-1]
			if prefix == "" {
				return ""
			}
		}
	}
	return prefix
}
