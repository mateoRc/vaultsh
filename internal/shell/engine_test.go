package shell

import "testing"

func TestEngineExecute(t *testing.T) {
	tests := []struct {
		name string
		line string
		want Result
	}{
		{
			name: "help",
			line: "help",
			want: Result{Output: "Available commands:\n  help", ExitCode: 0},
		},
		{
			name: "unknown command",
			line: "pwd",
			want: Result{Output: "command not found: pwd", ExitCode: 127},
		},
	}

	engine := &Engine{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := engine.Execute(tt.line)
			if got != tt.want {
				t.Errorf("Execute(%q) = %#v, want %#v", tt.line, got, tt.want)
			}
		})
	}
}
