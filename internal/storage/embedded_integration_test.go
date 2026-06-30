package storage_test

import (
	"strings"
	"testing"

	"github.com/mateom/vaultsh/content"
	"github.com/mateom/vaultsh/internal/command"
	"github.com/mateom/vaultsh/internal/shell"
	"github.com/mateom/vaultsh/internal/storage"
)

func TestEmbeddedContentThroughShell(t *testing.T) {
	root, err := storage.Load(content.Files)
	if err != nil {
		t.Fatalf("Load(content.Files): %v", err)
	}
	engine := shell.NewWithRoot(root)

	tests := []struct {
		name        string
		command     string
		wantOutput  string
		wantContain []string
	}{
		{
			name:       "root layout",
			command:    "ls",
			wantOutput: "cv/\ndocs/\nprojects/",
		},
		{
			name:       "hidden message of the day",
			command:    "ls -a",
			wantOutput: ".motd\ncv/\ndocs/\nprojects/",
		},
		{
			name:       "read message of the day",
			command:    "cat .motd",
			wantOutput: "Welcome to Vaultsh.\n",
		},
		{
			name:       "experience layout",
			command:    "ls cv/experience",
			wantOutput: "a1.txt\narisglobal.txt\nintellexi.txt\nreversinglabs.txt",
		},
		{
			name:    "cat embedded file",
			command: "cat cv/experience/reversinglabs.txt",
			wantContain: []string{
				"company: ReversingLabs",
				"responsibility: mentoring",
			},
		},
		{
			name:        "grep embedded file",
			command:     "grep '^technology:' cv/experience/reversinglabs.txt",
			wantContain: []string{"technology: Python", "technology: Docker"},
		},
		{
			name:    "tree embedded content",
			command: "tree",
			wantContain: []string{
				"about.txt",
				"interests.txt",
				"cv",
				"docs",
				"projects",
			},
		},
		{
			name:        "tree all embedded content",
			command:     "tree -a",
			wantContain: []string{".motd"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := engine.Execute(tt.command)
			if result.ExitCode != command.ExitSuccess {
				t.Fatalf("Execute(%q) exit code = %d, output = %q", tt.command, result.ExitCode, result.Output)
			}
			if tt.wantOutput != "" && result.Output != tt.wantOutput {
				t.Errorf("Execute(%q) output = %q, want %q", tt.command, result.Output, tt.wantOutput)
			}
			for _, expected := range tt.wantContain {
				if !strings.Contains(result.Output, expected) {
					t.Errorf("Execute(%q) output does not contain %q", tt.command, expected)
				}
			}
			if tt.command == "tree" && strings.Contains(result.Output, ".motd") {
				t.Error("tree output includes hidden .motd without -a")
			}
		})
	}

	for _, mutation := range []string{"rm about.txt", "touch new.txt", "mkdir private"} {
		result := engine.Execute(mutation)
		if result.ExitCode != command.ExitNotFound {
			t.Errorf("Execute(%q) exit code = %d, want %d", mutation, result.ExitCode, command.ExitNotFound)
		}
	}

	result := engine.Execute("ls")
	if strings.Contains(result.Output, "education.txt") {
		t.Error("root layout contains removed education.txt")
	}
}
