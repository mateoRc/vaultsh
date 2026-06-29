package parser

import (
	"errors"
	"reflect"
	"testing"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		name string
		line string
		want []string
	}{
		{
			name: "words",
			line: "cat about.txt",
			want: []string{"cat", "about.txt"},
		},
		{
			name: "whitespace",
			line: "  cat\tabout.txt  ",
			want: []string{"cat", "about.txt"},
		},
		{
			name: "quotes",
			line: `cat "my file.txt" 'second file.txt'`,
			want: []string{"cat", "my file.txt", "second file.txt"},
		},
		{
			name: "empty quoted argument",
			line: `cat ""`,
			want: []string{"cat", ""},
		},
		{
			name: "escape",
			line: `cat my\ file.txt`,
			want: []string{"cat", "my file.txt"},
		},
		{
			name: "escaped quote",
			line: `cat "my \"file\".txt"`,
			want: []string{"cat", `my "file".txt`},
		},
		{
			name: "pipe boundary",
			line: "cat about.txt|grep role",
			want: []string{"cat", "about.txt", "|", "grep", "role"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Tokenize(tt.line)
			if err != nil {
				t.Fatalf("Tokenize(): %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tokenize() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTokenizeErrors(t *testing.T) {
	tests := []struct {
		name string
		line string
		want error
	}{
		{
			name: "unterminated quote",
			line: `cat "about.txt`,
			want: ErrUnterminatedQuote,
		},
		{
			name: "trailing escape",
			line: `cat about.txt\`,
			want: ErrTrailingEscape,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Tokenize(tt.line)
			if !errors.Is(err, tt.want) {
				t.Errorf("Tokenize() error = %v, want %v", err, tt.want)
			}
		})
	}
}
