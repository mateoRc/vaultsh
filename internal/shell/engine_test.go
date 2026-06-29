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
					"\n  grep - Filter lines by a regular expression" +
					"\n  head - Print the first lines" +
					"\n  help - List available commands" +
					"\n  history - List commands from this session" +
					"\n  ls - List directory contents" +
					"\n  pwd - Print the current directory" +
					"\n  sort - Sort lines" +
					"\n  tail - Print the last lines" +
					"\n  tree - Print a directory tree" +
					"\n  wc - Count lines, words and bytes",
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

func TestEngineListOptions(t *testing.T) {
	root := filesystem.NewDirectory("")
	for _, node := range []filesystem.Node{
		filesystem.NewFile(".secret", "hidden"),
		filesystem.NewFile("about.txt", "hello"),
		filesystem.NewDirectory("docs"),
	} {
		if err := root.Add(node); err != nil {
			t.Fatalf("Add(%s): %v", node.Name(), err)
		}
	}
	engine := NewWithRoot(root)

	if result := engine.Execute("ls"); result.Output != "about.txt\ndocs/" {
		t.Errorf("ls output = %q", result.Output)
	}

	result := engine.Execute("ls -la")
	want := "-r--r--r--        6 .secret\n" +
		"-r--r--r--        5 about.txt\n" +
		"dr-xr-xr-x        - docs/"
	if result.Output != want {
		t.Errorf("ls -la output = %q, want %q", result.Output, want)
	}
}

func TestEngineListRejectsTimestampSortWithoutMetadata(t *testing.T) {
	result := New().Execute("ls -lt")

	if result.Output != "ls: option -t requires file timestamps" {
		t.Errorf("ls -lt output = %q", result.Output)
	}
	if result.ExitCode != command.ExitUnsupported {
		t.Errorf(
			"ls -lt exit code = %d, want %d",
			result.ExitCode,
			command.ExitUnsupported,
		)
	}
}

func TestEngineListRecursive(t *testing.T) {
	root := filesystem.NewDirectory("")
	docs := filesystem.NewDirectory("docs")
	nested := filesystem.NewDirectory("nested")
	if err := root.Add(filesystem.NewFile("about.txt", "hello")); err != nil {
		t.Fatalf("Add(about.txt): %v", err)
	}
	if err := root.Add(docs); err != nil {
		t.Fatalf("Add(docs): %v", err)
	}
	if err := docs.Add(filesystem.NewFile("readme.txt", "hello")); err != nil {
		t.Fatalf("Add(readme.txt): %v", err)
	}
	if err := docs.Add(nested); err != nil {
		t.Fatalf("Add(nested): %v", err)
	}
	if err := nested.Add(filesystem.NewFile("deep.txt", "hello")); err != nil {
		t.Fatalf("Add(deep.txt): %v", err)
	}

	result := NewWithRoot(root).Execute("ls -R")
	want := "/:\nabout.txt\ndocs/\n\n" +
		"/docs:\nnested/\nreadme.txt\n\n" +
		"/docs/nested:\ndeep.txt"

	if result.Output != want {
		t.Errorf("ls -R output = %q, want %q", result.Output, want)
	}
	if result.ExitCode != command.ExitSuccess {
		t.Errorf("exit code = %d, want %d", result.ExitCode, command.ExitSuccess)
	}
}

func TestEngineCommandHelp(t *testing.T) {
	result := New().Execute("help cat")

	if result.Output != "Usage: cat [-n] [file]\nPrint file contents" {
		t.Errorf(
			"help output = %q, want %q",
			result.Output,
			"Usage: cat [-n] [file]\nPrint file contents",
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

func TestEngineCommandUsage(t *testing.T) {
	engine := New()

	tests := []struct {
		line string
		want string
	}{
		{line: "cd one two", want: "usage: cd [directory]"},
		{line: "cat", want: "usage: cat [-n] [file]"},
		{line: "tree one two", want: "usage: tree [-L depth] [path]"},
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

func TestEngineCatLineNumbers(t *testing.T) {
	root := filesystem.NewDirectory("")
	if err := root.Add(filesystem.NewFile("about.txt", "first\nsecond\n")); err != nil {
		t.Fatalf("Add(about.txt): %v", err)
	}
	engine := NewWithRoot(root)

	tests := []string{
		"cat -n about.txt",
		"cat about.txt | cat -n",
	}
	want := "     1\tfirst\n     2\tsecond"

	for _, line := range tests {
		t.Run(line, func(t *testing.T) {
			result := engine.Execute(line)
			if result.Output != want {
				t.Errorf("output = %q, want %q", result.Output, want)
			}
			if result.ExitCode != command.ExitSuccess {
				t.Errorf(
					"exit code = %d, want %d",
					result.ExitCode,
					command.ExitSuccess,
				)
			}
		})
	}
}

func TestEnginePipeline(t *testing.T) {
	root := filesystem.NewDirectory("")
	if err := root.Add(filesystem.NewFile("about.txt", "hello")); err != nil {
		t.Fatalf("Add(about.txt): %v", err)
	}

	result := NewWithRoot(root).Execute("cat about.txt | cat")

	if result.Output != "hello" {
		t.Errorf("pipeline output = %q, want %q", result.Output, "hello")
	}
	if result.ExitCode != command.ExitSuccess {
		t.Errorf(
			"pipeline exit code = %d, want %d",
			result.ExitCode,
			command.ExitSuccess,
		)
	}
}

func TestEnginePipelineStopsOnFailure(t *testing.T) {
	root := filesystem.NewDirectory("")
	if err := root.Add(filesystem.NewFile("about.txt", "hello")); err != nil {
		t.Fatalf("Add(about.txt): %v", err)
	}

	result := NewWithRoot(root).Execute("cat missing.txt | cat")

	if result.ExitCode != command.ExitFailure {
		t.Errorf(
			"pipeline exit code = %d, want %d",
			result.ExitCode,
			command.ExitFailure,
		)
	}
}

func TestEngineGrepPipeline(t *testing.T) {
	root := filesystem.NewDirectory("")
	content := "language: Python\nlanguage: Go\nbackend: Flask\n"
	if err := root.Add(filesystem.NewFile("skills.txt", content)); err != nil {
		t.Fatalf("Add(skills.txt): %v", err)
	}

	result := NewWithRoot(root).Execute(`cat skills.txt | grep "^language:"`)

	want := "language: Python\nlanguage: Go"
	if result.Output != want {
		t.Errorf("grep output = %q, want %q", result.Output, want)
	}
	if result.ExitCode != command.ExitSuccess {
		t.Errorf("grep exit code = %d, want %d", result.ExitCode, command.ExitSuccess)
	}
}

func TestEngineGrepOptions(t *testing.T) {
	root := filesystem.NewDirectory("")
	content := "language: Python\nlanguage: Go\nbackend: Flask\n"
	if err := root.Add(filesystem.NewFile("skills.txt", content)); err != nil {
		t.Fatalf("Add(skills.txt): %v", err)
	}

	result := NewWithRoot(root).Execute("grep -in python skills.txt")

	if result.Output != "1:language: Python" {
		t.Errorf("grep output = %q, want %q", result.Output, "1:language: Python")
	}
	if result.ExitCode != command.ExitSuccess {
		t.Errorf("grep exit code = %d, want %d", result.ExitCode, command.ExitSuccess)
	}
}

func TestEngineGrepNoMatch(t *testing.T) {
	root := filesystem.NewDirectory("")
	if err := root.Add(filesystem.NewFile("skills.txt", "language: Go\n")); err != nil {
		t.Fatalf("Add(skills.txt): %v", err)
	}

	result := NewWithRoot(root).Execute("grep Python skills.txt")

	if result.Output != "" {
		t.Errorf("grep output = %q, want empty", result.Output)
	}
	if result.ExitCode != command.ExitFailure {
		t.Errorf("grep exit code = %d, want %d", result.ExitCode, command.ExitFailure)
	}
}

func TestEngineLineFilters(t *testing.T) {
	root := filesystem.NewDirectory("")
	if err := root.Add(filesystem.NewFile("values.txt", "c\nb\na\n")); err != nil {
		t.Fatalf("Add(values.txt): %v", err)
	}
	engine := NewWithRoot(root)

	tests := []struct {
		line string
		want string
	}{
		{line: "head -n 2 values.txt", want: "c\nb"},
		{line: "tail -n 2 values.txt", want: "b\na"},
		{line: "sort values.txt", want: "a\nb\nc"},
		{line: "sort -r values.txt", want: "c\nb\na"},
		{line: "wc -l values.txt", want: "3"},
		{line: "wc values.txt", want: "3 3 6"},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			result := engine.Execute(tt.line)
			if result.Output != tt.want {
				t.Errorf("output = %q, want %q", result.Output, tt.want)
			}
			if result.ExitCode != command.ExitSuccess {
				t.Errorf(
					"exit code = %d, want %d",
					result.ExitCode,
					command.ExitSuccess,
				)
			}
		})
	}
}

func TestEngineMultiStagePipeline(t *testing.T) {
	root := filesystem.NewDirectory("")
	content := "language: Python\nbackend: Flask\nlanguage: Go\nlanguage: Java\n"
	if err := root.Add(filesystem.NewFile("skills.txt", content)); err != nil {
		t.Fatalf("Add(skills.txt): %v", err)
	}

	result := NewWithRoot(root).Execute(
		`cat skills.txt | grep "^language:" | sort | head -n 2`,
	)

	want := "language: Go\nlanguage: Java"
	if result.Output != want {
		t.Errorf("pipeline output = %q, want %q", result.Output, want)
	}
	if result.ExitCode != command.ExitSuccess {
		t.Errorf("exit code = %d, want %d", result.ExitCode, command.ExitSuccess)
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

func TestEngineTreeDepth(t *testing.T) {
	root := filesystem.NewDirectory("")
	docs := filesystem.NewDirectory("docs")
	nested := filesystem.NewDirectory("nested")
	if err := root.Add(docs); err != nil {
		t.Fatalf("Add(docs): %v", err)
	}
	if err := docs.Add(nested); err != nil {
		t.Fatalf("Add(nested): %v", err)
	}
	if err := nested.Add(filesystem.NewFile("readme.txt", "hello")); err != nil {
		t.Fatalf("Add(readme.txt): %v", err)
	}

	result := NewWithRoot(root).Execute("tree -L 1")
	want := ".\n└── docs"

	if result.Output != want {
		t.Errorf("tree output = %q, want %q", result.Output, want)
	}
	if result.ExitCode != command.ExitSuccess {
		t.Errorf("exit code = %d, want %d", result.ExitCode, command.ExitSuccess)
	}
}

func TestEngineTreeRejectsInvalidDepth(t *testing.T) {
	result := New().Execute("tree -L 0")

	if result.Output != "tree: invalid depth: 0" {
		t.Errorf("tree output = %q", result.Output)
	}
	if result.ExitCode != command.ExitUsage {
		t.Errorf("exit code = %d, want %d", result.ExitCode, command.ExitUsage)
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

func TestEngineHireEasterEgg(t *testing.T) {
	engine := New()

	result := engine.Execute("hire mateo")

	want := "hire: permission denied\nhint: try sudo hire mateo -s <salary>"
	if result.Output != want {
		t.Errorf("hire output = %q, want %q", result.Output, want)
	}
	if result.ExitCode != command.ExitFailure {
		t.Errorf("exit code = %d, want %d", result.ExitCode, command.ExitFailure)
	}
}

func TestEngineSudoHireEasterEgg(t *testing.T) {
	result := New().Execute("sudo hire mateo -s 100000")

	want := "sudo: access granted\n" +
		"salary offered: 100000.00\n" +
		"counter-offer: 150000.00\n" +
		"accept counter-offer? [Y/y]"
	if result.Output != want {
		t.Errorf("sudo hire output = %q, want %q", result.Output, want)
	}
	if result.ExitCode != command.ExitSuccess {
		t.Errorf("exit code = %d, want %d", result.ExitCode, command.ExitSuccess)
	}
}

func TestEngineAcceptsCounterOfferOnce(t *testing.T) {
	engine := New()
	if result := engine.Execute("sudo hire mateo -s 100000"); result.ExitCode != command.ExitSuccess {
		t.Fatalf("sudo hire failed: %s", result.Output)
	}

	result := engine.Execute("y")
	want := "counter-offer accepted: 150000.00\n" +
		"welcome aboard. paperwork has entered the chat."
	if result.Output != want {
		t.Errorf("accept output = %q, want %q", result.Output, want)
	}
	if result.ExitCode != command.ExitSuccess {
		t.Errorf("exit code = %d, want %d", result.ExitCode, command.ExitSuccess)
	}

	result = engine.Execute("Y")
	if result.Output != "Y: no pending counter-offer" {
		t.Errorf("second accept output = %q", result.Output)
	}
	if result.ExitCode != command.ExitFailure {
		t.Errorf("exit code = %d, want %d", result.ExitCode, command.ExitFailure)
	}
}

func TestEngineCounterOfferIsSessionScoped(t *testing.T) {
	first := New()
	second := New()
	first.Execute("sudo hire mateo -s 100000")

	result := second.Execute("y")

	if result.Output != "y: no pending counter-offer" {
		t.Errorf("second session output = %q", result.Output)
	}
	if result.ExitCode != command.ExitFailure {
		t.Errorf("exit code = %d, want %d", result.ExitCode, command.ExitFailure)
	}
}

func TestEngineSudoHireRequiresValidSalary(t *testing.T) {
	tests := []string{
		"sudo hire mateo",
		"sudo hire mateo -s nope",
		"sudo hire mateo -s 0",
	}

	for _, line := range tests {
		t.Run(line, func(t *testing.T) {
			result := New().Execute(line)
			if result.ExitCode != command.ExitUsage {
				t.Errorf(
					"exit code = %d, want %d",
					result.ExitCode,
					command.ExitUsage,
				)
			}
		})
	}
}
