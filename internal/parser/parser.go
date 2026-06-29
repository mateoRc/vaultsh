package parser

import "errors"

var ErrUnexpectedPipe = errors.New("unexpected pipe")

func Parse(tokens []Token) (SyntaxTree, error) {
	if len(tokens) == 0 {
		return SyntaxTree{}, nil
	}

	commands := make([]CommandNode, 1)
	for _, token := range tokens {
		if token.Kind == TokenPipe {
			if commands[len(commands)-1].Name == "" {
				return SyntaxTree{}, ErrUnexpectedPipe
			}
			commands = append(commands, CommandNode{})
			continue
		}

		current := &commands[len(commands)-1]
		if current.Name == "" {
			current.Name = token.Value
			continue
		}
		current.Args = append(current.Args, token.Value)
	}

	if commands[len(commands)-1].Name == "" {
		return SyntaxTree{}, ErrUnexpectedPipe
	}

	return SyntaxTree{Pipeline: commands}, nil
}
