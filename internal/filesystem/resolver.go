package filesystem

import (
	"errors"
	"fmt"
	"path"
	"strings"
)

var (
	ErrNodeNotFound = errors.New("node not found")
	ErrNotDirectory = errors.New("not a directory")
)

type Resolver struct {
	root *Directory
}

func NewResolver(root *Directory) *Resolver {
	return &Resolver{root: root}
}

func (r *Resolver) Resolve(workingDirectory, target string) (Node, string, error) {
	resolvedPath := resolvePath(workingDirectory, target)
	if resolvedPath == "/" {
		return r.root, resolvedPath, nil
	}

	var current Node = r.root
	for _, name := range strings.Split(strings.TrimPrefix(resolvedPath, "/"), "/") {
		directory, ok := current.(*Directory)
		if !ok {
			return nil, "", fmt.Errorf("%w: %s", ErrNotDirectory, current.Name())
		}

		next, found := directory.Find(name)
		if !found {
			return nil, "", fmt.Errorf("%w: %s", ErrNodeNotFound, name)
		}
		current = next
	}

	return current, resolvedPath, nil
}

func resolvePath(workingDirectory, target string) string {
	if strings.HasPrefix(target, "/") {
		return path.Clean(target)
	}

	return path.Clean(path.Join("/", workingDirectory, target))
}
