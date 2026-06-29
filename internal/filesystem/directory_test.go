package filesystem

import (
	"errors"
	"testing"
)

type stubNode struct {
	name string
	kind Kind
}

func (n stubNode) Name() string {
	return n.name
}

func (n stubNode) Kind() Kind {
	return n.kind
}

func TestDirectoryImplementsNode(t *testing.T) {
	var node Node = NewDirectory("docs")

	if node.Name() != "docs" {
		t.Errorf("Name() = %q, want %q", node.Name(), "docs")
	}
	if node.Kind() != KindDirectory {
		t.Errorf("Kind() = %q, want %q", node.Kind(), KindDirectory)
	}
}

func TestDirectoryAddAndFind(t *testing.T) {
	directory := NewDirectory("root")
	expected := stubNode{name: "readme", kind: KindFile}

	if err := directory.Add(expected); err != nil {
		t.Fatalf("Add(): %v", err)
	}

	actual, found := directory.Find("readme")
	if !found {
		t.Fatal("Find(readme) did not find added node")
	}
	if actual != expected {
		t.Errorf("Find(readme) = %#v, want %#v", actual, expected)
	}
}

func TestDirectoryRejectsDuplicateName(t *testing.T) {
	directory := NewDirectory("root")
	if err := directory.Add(stubNode{name: "docs"}); err != nil {
		t.Fatalf("first Add(): %v", err)
	}

	err := directory.Add(stubNode{name: "docs"})

	if !errors.Is(err, ErrNodeExists) {
		t.Errorf("second Add() error = %v, want %v", err, ErrNodeExists)
	}
}

func TestDirectoryChildrenAreSorted(t *testing.T) {
	directory := NewDirectory("root")
	for _, name := range []string{"zebra", "alpha"} {
		if err := directory.Add(stubNode{name: name}); err != nil {
			t.Fatalf("Add(%s): %v", name, err)
		}
	}

	children := directory.Children()

	if len(children) != 2 {
		t.Fatalf("len(Children()) = %d, want 2", len(children))
	}
	if children[0].Name() != "alpha" || children[1].Name() != "zebra" {
		t.Errorf(
			"Children() names = [%q, %q], want [alpha, zebra]",
			children[0].Name(),
			children[1].Name(),
		)
	}
}
