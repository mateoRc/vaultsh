package content

import (
	"io/fs"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"
	"unicode/utf8"
)

var contentLinePattern = regexp.MustCompile(`^[a-z][a-z0-9_]*: .+$`)

func TestEmbeddedContentLayout(t *testing.T) {
	want := []string{
		".motd",
		"cv/about.txt",
		"cv/experience/a1.txt",
		"cv/experience/arisglobal.txt",
		"cv/experience/intellexi.txt",
		"cv/experience/reversinglabs.txt",
		"cv/interests.txt",
		"cv/skills.txt",
		"docs/api.md",
		"docs/commands.md",
		"docs/content.md",
		"docs/roadmap.md",
		"projects/vaultsh.txt",
	}

	var got []string
	err := fs.WalkDir(Files, ".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path != "." && !entry.IsDir() {
			got = append(got, path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("WalkDir(Files): %v", err)
	}
	sort.Strings(got)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("embedded content paths = %v, want %v", got, want)
	}
}

func TestEmbeddedPlainTextFormat(t *testing.T) {
	err := fs.WalkDir(Files, ".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || strings.HasSuffix(path, ".md") {
			return nil
		}

		data, err := fs.ReadFile(Files, path)
		if err != nil {
			t.Errorf("ReadFile(%q): %v", path, err)
			return nil
		}
		if !utf8.Valid(data) {
			t.Errorf("%s is not valid UTF-8", path)
			return nil
		}
		if strings.Contains(string(data), "\r") {
			t.Errorf("%s uses CRLF; content files must use LF", path)
		}
		if path == ".motd" {
			if string(data) != "Welcome to Vaultsh.\n" {
				t.Errorf(".motd = %q, want %q", string(data), "Welcome to Vaultsh.\n")
			}
			return nil
		}

		for number, line := range strings.Split(string(data), "\n") {
			if line != "" && !contentLinePattern.MatchString(line) {
				t.Errorf("%s:%d is not a key-value line: %q", path, number+1, line)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("WalkDir(Files): %v", err)
	}
}
