package shell

import (
	"testing"

	"github.com/mateom/vaultsh/internal/command"
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
					"\n  help - List available commands",
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
			name: "unknown command",
			line: "pwd",
			want: command.Result{
				Output:   "command not found: pwd",
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
