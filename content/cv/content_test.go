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
		"about.txt",
		"experience/a1.txt",
		"experience/arisglobal.txt",
		"experience/intellexi.txt",
		"experience/reversinglabs.txt",
		"interests.txt",
		"projects/vaultsh.txt",
		"skills.txt",
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

func TestEmbeddedContentFormat(t *testing.T) {
	err := fs.WalkDir(Files, ".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
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

		for number, line := range strings.Split(string(data), "\n") {
			if line == "" {
				continue
			}
			if strings.HasPrefix(line, "#") {
				t.Errorf("%s:%d uses a Markdown heading", path, number+1)
				continue
			}
			if !contentLinePattern.MatchString(line) {
				t.Errorf("%s:%d is not a key-value line: %q", path, number+1, line)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("WalkDir(Files): %v", err)
	}
}
