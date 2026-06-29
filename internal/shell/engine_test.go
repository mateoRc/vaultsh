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
					"\n  cat - Print file contents" +
					"\n  cd - Change the current directory" +
					"\n  clear - Clear the terminal" +
					"\n  help - List available commands" +
					"\n  history - List commands from this session" +
					"\n  ls - List directory contents" +
					"\n  pwd - Print the current directory" +
					"\n  tree - Print a directory tree",
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

	if result.Output != "about.txt\ndocs/" {
		t.Errorf("ls output = %q, want %q", result.Output, "about.txt\ndocs/")
	}
	if result.ExitCode != 0 {
		t.Errorf("ls exit code = %d, want 0", result.ExitCode)
	}
}

func TestEngineListPath(t *testing.T) {
	root := filesystem.NewDirectory("")
	docs := filesystem.NewDirectory("docs")
	if err := root.Add(docs); err != nil {
		t.Fatalf("Add(docs): %v", err)
	}
	if err := docs.Add(filesystem.NewFile("readme.txt", "hello")); err != nil {
		t.Fatalf("Add(readme.txt): %v", err)
	}

	result := NewWithRoot(root).Execute("ls /docs")

	if result.Output != "readme.txt" {
		t.Errorf("ls output = %q, want %q", result.Output, "readme.txt")
	}
	if result.ExitCode != 0 {
		t.Errorf("ls exit code = %d, want 0", result.ExitCode)
	}
}

func TestEngineListFile(t *testing.T) {
	root := filesystem.NewDirectory("")
	if err := root.Add(filesystem.NewFile("about.txt", "hello")); err != nil {
		t.Fatalf("Add(about.txt): %v", err)
	}

	result := NewWithRoot(root).Execute("ls about.txt")

	if result.Output != "about.txt" {
		t.Errorf("ls output = %q, want %q", result.Output, "about.txt")
	}
	if result.ExitCode != command.ExitSuccess {
		t.Errorf("ls exit code = %d, want %d", result.ExitCode, command.ExitSuccess)
	}
}

func TestEngineCommandHelp(t *testing.T) {
	result := New().Execute("help cat")

	if result.Output != "Usage: cat <file>\nPrint file contents" {
		t.Errorf(
			"help output = %q, want %q",
			result.Output,
			"Usage: cat <file>\nPrint file contents",
		)
	}
	if result.ExitCode != 0 {
		t.Errorf("help exit code = %d, want 0", result.ExitCode)
	}
}

func TestEngineChangeDirectoryAndReadFile(t *testing.T) {
	root := filesystem.NewDirectory("")
	docs := filesystem.NewDirectory("docs")
	if err := root.Add(docs); err != nil {
		t.Fatalf("Add(docs): %v", err)
	}
	if err := docs.Add(filesystem.NewFile("readme.txt", "hello")); err != nil {
		t.Fatalf("Add(readme.txt): %v", err)
	}
	engine := NewWithRoot(root)

	if result := engine.Execute("cd docs"); result.ExitCode != 0 {
		t.Fatalf("cd exit code = %d, output = %q", result.ExitCode, result.Output)
	}
	if result := engine.Execute("pwd"); result.Output != "/docs" {
		t.Errorf("pwd output = %q, want %q", result.Output, "/docs")
	}
	if result := engine.Execute("cat readme.txt"); result.Output != "hello" {
		t.Errorf("cat output = %q, want %q", result.Output, "hello")
	}
}

func TestEngineQuotedArgument(t *testing.T) {
	root := filesystem.NewDirectory("")
	if err := root.Add(filesystem.NewFile("my file.txt", "hello")); err != nil {
		t.Fatalf("Add(my file.txt): %v", err)
	}

	result := NewWithRoot(root).Execute(`cat "my file.txt"`)

	if result.Output != "hello" {
		t.Errorf("cat output = %q, want %q", result.Output, "hello")
	}
	if result.ExitCode != 0 {
		t.Errorf("cat exit code = %d, want 0", result.ExitCode)
	}
}

func TestEngineSyntaxError(t *testing.T) {
	result := New().Execute(`cat "about.txt`)

	if result.Output != "syntax error: unterminated quote" {
		t.Errorf(
			"output = %q, want %q",
			result.Output,
			"syntax error: unterminated quote",
		)
	}
	if result.ExitCode != 2 {
		t.Errorf("exit code = %d, want 2", result.ExitCode)
	}
}

func TestEngineRejectsPipelineUntilSupported(t *testing.T) {
	result := New().Execute("cat about.txt | grep role")

	if result.Output != "pipelines are not supported yet" {
		t.Errorf(
			"output = %q, want %q",
			result.Output,
			"pipelines are not supported yet",
		)
	}
	if result.ExitCode != command.ExitUnsupported {
		t.Errorf("exit code = %d, want %d", result.ExitCode, command.ExitUnsupported)
	}
}

func TestEngineCommandUsage(t *testing.T) {
	engine := New()

	tests := []struct {
		line string
		want string
	}{
		{line: "cd one two", want: "usage: cd [directory]"},
		{line: "cat", want: "usage: cat <file>"},
		{line: "tree one two", want: "usage: tree [path]"},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			result := engine.Execute(tt.line)
			if result.ExitCode != 2 {
				t.Errorf("exit code = %d, want 2", result.ExitCode)
			}
			if result.Output != tt.want {
				t.Errorf("output = %q, want %q", result.Output, tt.want)
			}
		})
	}
}

func TestEngineCatRejectsDirectory(t *testing.T) {
	root := filesystem.NewDirectory("")
	if err := root.Add(filesystem.NewDirectory("docs")); err != nil {
		t.Fatalf("Add(docs): %v", err)
	}

	result := NewWithRoot(root).Execute("cat docs")

	if result.Output != "cat: docs: is a directory" {
		t.Errorf("cat output = %q, want %q", result.Output, "cat: docs: is a directory")
	}
	if result.ExitCode != 1 {
		t.Errorf("cat exit code = %d, want 1", result.ExitCode)
	}
}

func TestEngineTree(t *testing.T) {
	root := filesystem.NewDirectory("")
	docs := filesystem.NewDirectory("docs")
	if err := root.Add(filesystem.NewFile("about.txt", "hello")); err != nil {
		t.Fatalf("Add(about.txt): %v", err)
	}
	if err := root.Add(docs); err != nil {
		t.Fatalf("Add(docs): %v", err)
	}
	if err := docs.Add(filesystem.NewFile("readme.txt", "hello")); err != nil {
		t.Fatalf("Add(readme.txt): %v", err)
	}

	result := NewWithRoot(root).Execute("tree")
	want := ".\n├── about.txt\n└── docs\n    └── readme.txt"

	if result.Output != want {
		t.Errorf("tree output = %q, want %q", result.Output, want)
	}
	if result.ExitCode != 0 {
		t.Errorf("tree exit code = %d, want 0", result.ExitCode)
	}
}

func TestEnginesKeepIndependentWorkingDirectories(t *testing.T) {
	root := filesystem.NewDirectory("")
	if err := root.Add(filesystem.NewDirectory("docs")); err != nil {
		t.Fatalf("Add(docs): %v", err)
	}
	first := NewWithContext(NewExecutionContext(root))
	second := NewWithContext(NewExecutionContext(root))

	if result := first.Execute("cd docs"); result.ExitCode != 0 {
		t.Fatalf("cd exit code = %d, output = %q", result.ExitCode, result.Output)
	}

	if result := first.Execute("pwd"); result.Output != "/docs" {
		t.Errorf("first pwd output = %q, want /docs", result.Output)
	}
	if result := second.Execute("pwd"); result.Output != "/" {
		t.Errorf("second pwd output = %q, want /", result.Output)
	}
}

func TestEngineHistory(t *testing.T) {
	engine := New()
	engine.Execute("pwd")
	engine.Execute("about")

	result := engine.Execute("history")
	want := "1  pwd\n2  about\n3  history"

	if result.Output != want {
		t.Errorf("history output = %q, want %q", result.Output, want)
	}
	if result.ExitCode != 0 {
		t.Errorf("history exit code = %d, want 0", result.ExitCode)
	}
}
