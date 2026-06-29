package parser

import (
	"errors"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tokens := []Token{
		{Kind: TokenWord, Value: "cat"},
		{Kind: TokenWord, Value: "about.txt"},
		{Kind: TokenPipe, Value: "|"},
		{Kind: TokenWord, Value: "grep"},
		{Kind: TokenWord, Value: "role"},
	}

	got, err := Parse(tokens)
	if err != nil {
		t.Fatalf("Parse(): %v", err)
	}

	want := [][]string{
		{"cat", "about.txt"},
		{"grep", "role"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Parse() = %#v, want %#v", got, want)
	}
}

func TestParseEmptyInput(t *testing.T) {
	got, err := Parse(nil)

	if err != nil {
		t.Fatalf("Parse(): %v", err)
	}
	if got != nil {
		t.Errorf("Parse() = %#v, want nil", got)
	}
}

func TestParseRejectsUnexpectedPipe(t *testing.T) {
	tests := []struct {
		name   string
		tokens []Token
	}{
		{
			name: "leading",
			tokens: []Token{
				{Kind: TokenPipe, Value: "|"},
				{Kind: TokenWord, Value: "cat"},
			},
		},
		{
			name: "only pipe",
			tokens: []Token{
				{Kind: TokenPipe, Value: "|"},
			},
		},
		{
			name: "trailing",
			tokens: []Token{
				{Kind: TokenWord, Value: "cat"},
				{Kind: TokenPipe, Value: "|"},
			},
		},
		{
			name: "consecutive",
			tokens: []Token{
				{Kind: TokenWord, Value: "cat"},
				{Kind: TokenPipe, Value: "|"},
				{Kind: TokenPipe, Value: "|"},
				{Kind: TokenWord, Value: "grep"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.tokens)
			if !errors.Is(err, ErrUnexpectedPipe) {
				t.Errorf("Parse() error = %v, want %v", err, ErrUnexpectedPipe)
			}
		})
	}
}
