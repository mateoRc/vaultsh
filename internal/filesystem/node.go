// Package filesystem provides Vaultsh's virtual filesystem model.
package filesystem

// Kind identifies the type of a filesystem node.
type Kind string

const (
	KindFile      Kind = "file"
	KindDirectory Kind = "directory"
)

type Node interface {
	Name() string
	Kind() Kind
}
