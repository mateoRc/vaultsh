package parser

import (
	"reflect"
	"testing"
)

func TestLex(t *testing.T) {
	got, err := Lex(`cat "about file.txt"|grep role`)
	if err != nil {
		t.Fatalf("Lex(): %v", err)
	}

	want := []Token{
		{Kind: TokenWord, Value: "cat"},
		{Kind: TokenWord, Value: "about file.txt"},
		{Kind: TokenPipe, Value: "|"},
		{Kind: TokenWord, Value: "grep"},
		{Kind: TokenWord, Value: "role"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Lex() = %#v, want %#v", got, want)
	}
}

func TestLexPropagatesTokenizerError(t *testing.T) {
	_, err := Lex(`cat "about.txt`)

	if err != ErrUnterminatedQuote {
		t.Errorf("Lex() error = %v, want %v", err, ErrUnterminatedQuote)
	}
}
