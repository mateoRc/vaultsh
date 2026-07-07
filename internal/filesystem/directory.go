package filesystem

import (
	"errors"
	"sort"
	"time"
)

var ErrNodeExists = errors.New("node already exists")

type Directory struct {
	name     string
	modTime  time.Time
	children map[string]Node
}

func NewDirectory(name string) *Directory {
	return NewDirectoryWithModTime(name, time.Time{})
}

func NewDirectoryWithModTime(name string, modTime time.Time) *Directory {
	return &Directory{
		name:     name,
		modTime:  modTime,
		children: make(map[string]Node),
	}
}

func (d *Directory) Name() string {
	return d.name
}

func (*Directory) Kind() Kind {
	return KindDirectory
}

func (d *Directory) ModTime() time.Time {
	return d.modTime
}

func (d *Directory) Add(node Node) error {
	if _, exists := d.children[node.Name()]; exists {
		return ErrNodeExists
	}

	d.children[node.Name()] = node
	return nil
}

func (d *Directory) Find(name string) (Node, bool) {
	node, found := d.children[name]
	return node, found
}

func (d *Directory) Children() []Node {
	children := make([]Node, 0, len(d.children))
	for _, node := range d.children {
		children = append(children, node)
	}

	sort.Slice(children, func(i, j int) bool {
		return children[i].Name() < children[j].Name()
	})

	return children
}
