package shell

import (
	"testing"

	"github.com/mateom/vaultsh/internal/command"
	"github.com/mateom/vaultsh/internal/filesystem"
)

func TestEngineExecute(t *testing.T) {
	tests := []struct {
		name string
		line string
		want command.Result
	}{
		{
			name: "help",
			line: "help",
			want: command.Result{
				Output: "Available commands:" +
					"\n  about - Describe Vaultsh" +
					"\n  clear - Clear the terminal" +
					"\n  help - List available commands" +
					"\n  ls - List directory contents" +
					"\n  pwd - Print the current directory",
				ExitCode: 0,
			},
		},
		{
			name: "about",
			line: "about",
			want: command.Result{
				Output:   "Vaultsh is a read-only virtual shell engine.",
				ExitCode: 0,
			},
		},
		{
			name: "clear",
			line: "clear",
			want: command.Result{
				ExitCode: 0,
				Action:   command.ActionClear,
			},
		},
		{
			name: "pwd",
			line: "pwd",
			want: command.Result{
				Output:   "/",
				ExitCode: 0,
			},
		},
		{
			name: "ls empty directory",
			line: "ls",
			want: command.Result{
				Output:   "",
				ExitCode: 0,
			},
		},
		{
			name: "unknown command",
			line: "missing",
			want: command.Result{
				Output:   "command not found: missing",
				ExitCode: 127,
			},
		},
	}

	engine := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := engine.Execute(tt.line)
			if got != tt.want {
				t.Errorf("Execute(%q) = %#v, want %#v", tt.line, got, tt.want)
			}
		})
	}
}

func TestEngineListDirectory(t *testing.T) {
	root := filesystem.NewDirectory("")
	if err := root.Add(filesystem.NewDirectory("docs")); err != nil {
		t.Fatalf("Add(docs): %v", err)
	}
	if err := root.Add(filesystem.NewFile("about.txt", "hello")); err != nil {
		t.Fatalf("Add(about.txt): %v", err)
	}

	result := NewWithRoot(root).Execute("ls")

	if result.Output != "about.txt\ndocs" {
		t.Errorf("ls output = %q, want %q", result.Output, "about.txt\ndocs")
	}
	if result.ExitCode != 0 {
		t.Errorf("ls exit code = %d, want 0", result.ExitCode)
	}
}
