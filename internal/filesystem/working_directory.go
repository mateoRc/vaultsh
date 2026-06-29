package filesystem

import "fmt"

type WorkingDirectory struct {
	resolver  *Resolver
	directory *Directory
	path      string
}

func NewWorkingDirectory(root *Directory) *WorkingDirectory {
	return &WorkingDirectory{
		resolver:  NewResolver(root),
		directory: root,
		path:      "/",
	}
}

func (w *WorkingDirectory) Path() string {
	return w.path
}

func (w *WorkingDirectory) Directory() *Directory {
	return w.directory
}

func (w *WorkingDirectory) Resolve(target string) (Node, string, error) {
	return w.resolver.Resolve(w.path, target)
}

func (w *WorkingDirectory) Change(target string) error {
	node, resolvedPath, err := w.Resolve(target)
	if err != nil {
		return err
	}

	directory, ok := node.(*Directory)
	if !ok {
		return fmt.Errorf("%w: %s", ErrNotDirectory, resolvedPath)
	}

	w.directory = directory
	w.path = resolvedPath
	return nil
}
