package filesystem

type Kind string

const (
	KindFile      Kind = "file"
	KindDirectory Kind = "directory"
)

type Node interface {
	Name() string
	Kind() Kind
}
