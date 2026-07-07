package command

import (
	"strings"
	"testing"
	"time"

	"github.com/mateom/vaultsh/internal/filesystem"
)

func TestLsReverseOrder(t *testing.T) {
	root := filesystem.NewDirectory("")
	mustAddNode(t, root, filesystem.NewFile("alpha.md", "alpha"))
	mustAddNode(t, root, filesystem.NewFile("zebra.md", "zebra"))

	result := NewLs(filesystem.NewWorkingDirectory(root)).Execute(
		[]string{"-r"},
		Input{},
	)

	if result.ExitCode != ExitSuccess {
		t.Fatalf("ExitCode = %d, want %d; output=%q", result.ExitCode, ExitSuccess, result.Output)
	}
	if result.Output != "zebra.md\nalpha.md" {
		t.Errorf("Output = %q, want reverse alphabetical listing", result.Output)
	}
}

func TestLsLongReverseOrder(t *testing.T) {
	root := filesystem.NewDirectory("")
	mustAddNode(t, root, filesystem.NewFile("alpha.md", "alpha"))
	mustAddNode(t, root, filesystem.NewFile("zebra.md", "zebra"))

	result := NewLs(filesystem.NewWorkingDirectory(root)).Execute(
		[]string{"-lr"},
		Input{},
	)

	if result.ExitCode != ExitSuccess {
		t.Fatalf("ExitCode = %d, want %d; output=%q", result.ExitCode, ExitSuccess, result.Output)
	}
	lines := strings.Split(result.Output, "\n")
	if len(lines) != 2 {
		t.Fatalf("lines = %d, want 2; output=%q", len(lines), result.Output)
	}
	if !strings.HasSuffix(lines[0], "zebra.md") || !strings.HasSuffix(lines[1], "alpha.md") {
		t.Errorf("Output = %q, want long reverse alphabetical listing", result.Output)
	}
}

func TestLsTimeSort(t *testing.T) {
	root := filesystem.NewDirectory("")
	older := time.Date(2026, 7, 6, 12, 0, 0, 0, time.UTC)
	newer := time.Date(2026, 7, 7, 12, 0, 0, 0, time.UTC)
	mustAddNode(t, root, filesystem.NewFileWithModTime("older.md", "older", older))
	mustAddNode(t, root, filesystem.NewFileWithModTime("newer.md", "newer", newer))

	result := NewLs(filesystem.NewWorkingDirectory(root)).Execute(
		[]string{"-t"},
		Input{},
	)

	if result.ExitCode != ExitSuccess {
		t.Fatalf("ExitCode = %d, want %d; output=%q", result.ExitCode, ExitSuccess, result.Output)
	}
	if result.Output != "newer.md\nolder.md" {
		t.Errorf("Output = %q, want newest file first", result.Output)
	}
}

func TestLsTimeSortCanBeReversed(t *testing.T) {
	root := filesystem.NewDirectory("")
	older := time.Date(2026, 7, 6, 12, 0, 0, 0, time.UTC)
	newer := time.Date(2026, 7, 7, 12, 0, 0, 0, time.UTC)
	mustAddNode(t, root, filesystem.NewFileWithModTime("older.md", "older", older))
	mustAddNode(t, root, filesystem.NewFileWithModTime("newer.md", "newer", newer))

	result := NewLs(filesystem.NewWorkingDirectory(root)).Execute(
		[]string{"-tr"},
		Input{},
	)

	if result.ExitCode != ExitSuccess {
		t.Fatalf("ExitCode = %d, want %d; output=%q", result.ExitCode, ExitSuccess, result.Output)
	}
	if result.Output != "older.md\nnewer.md" {
		t.Errorf("Output = %q, want oldest file first", result.Output)
	}
}

func mustAddNode(t *testing.T, directory *filesystem.Directory, node filesystem.Node) {
	t.Helper()
	if err := directory.Add(node); err != nil {
		t.Fatalf("Add(%s): %v", node.Name(), err)
	}
}
