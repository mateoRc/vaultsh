package parser

type TokenKind uint8

const (
	TokenWord TokenKind = iota
	TokenPipe
)

type Token struct {
	Kind  TokenKind
	Value string
}

func Lex(line string) ([]Token, error) {
	values, err := Tokenize(line)
	if err != nil {
		return nil, err
	}

	tokens := make([]Token, 0, len(values))
	for _, value := range values {
		kind := TokenWord
		if value == pipeToken {
			kind = TokenPipe
		}
		tokens = append(tokens, Token{
			Kind:  kind,
			Value: value,
		})
	}

	return tokens, nil
}
