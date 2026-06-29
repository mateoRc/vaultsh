package parser

import (
	"errors"
	"strings"
	"unicode"
)

const (
	pipeToken       = "|"
	pipeCharacter   = '|'
	escapeCharacter = '\\'
	singleQuote     = '\''
	doubleQuote     = '"'
)

var (
	ErrUnterminatedQuote = errors.New("unterminated quote")
	ErrTrailingEscape    = errors.New("trailing escape")
)

func Tokenize(line string) ([]string, error) {
	var tokens []string
	var current strings.Builder
	var quote rune
	escaped := false
	started := false

	flush := func() {
		if !started {
			return
		}
		tokens = append(tokens, current.String())
		current.Reset()
		started = false
	}

	for _, character := range line {
		if escaped {
			current.WriteRune(character)
			escaped = false
			started = true
			continue
		}

		if character == escapeCharacter && quote != singleQuote {
			escaped = true
			started = true
			continue
		}

		if quote != 0 {
			if character == quote {
				quote = 0
			} else {
				current.WriteRune(character)
			}
			continue
		}

		switch {
		case character == singleQuote || character == doubleQuote:
			quote = character
			started = true
		case unicode.IsSpace(character):
			flush()
		case character == pipeCharacter:
			flush()
			tokens = append(tokens, pipeToken)
		default:
			current.WriteRune(character)
			started = true
		}
	}

	if escaped {
		return nil, ErrTrailingEscape
	}
	if quote != 0 {
		return nil, ErrUnterminatedQuote
	}

	flush()
	return tokens, nil
}
