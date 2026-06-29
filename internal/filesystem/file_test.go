package filesystem

import "testing"

func TestFileImplementsNode(t *testing.T) {
	var node Node = NewFile("readme.txt", "hello")

	if node.Name() != "readme.txt" {
		t.Errorf("Name() = %q, want %q", node.Name(), "readme.txt")
	}
	if node.Kind() != KindFile {
		t.Errorf("Kind() = %q, want %q", node.Kind(), KindFile)
	}
}

func TestFileContent(t *testing.T) {
	file := NewFile("readme.txt", "hello")

	if file.Content() != "hello" {
		t.Errorf("Content() = %q, want %q", file.Content(), "hello")
	}
}
