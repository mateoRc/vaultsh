package parser

type CommandNode struct {
	Name string
	Args []string
}

type SyntaxTree struct {
	Pipeline []CommandNode
}
