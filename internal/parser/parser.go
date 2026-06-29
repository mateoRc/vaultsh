package parser

import "errors"

var ErrUnexpectedPipe = errors.New("unexpected pipe")

func Parse(tokens []Token) ([][]string, error) {
	if len(tokens) == 0 {
		return nil, nil
	}

	commands := make([][]string, 1)
	for _, token := range tokens {
		if token.Kind == TokenPipe {
			if len(commands[len(commands)-1]) == 0 {
				return nil, ErrUnexpectedPipe
			}
			commands = append(commands, nil)
			continue
		}

		commands[len(commands)-1] = append(commands[len(commands)-1], token.Value)
	}

	if len(commands[len(commands)-1]) == 0 {
		return nil, ErrUnexpectedPipe
	}

	return commands, nil
}
