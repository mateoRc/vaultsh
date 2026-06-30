package content

import "context"

// Provider loads the domain content exposed by the virtual filesystem.
// Implementations own persistence details; callers work only with this catalog.
type Provider interface {
	Load(context.Context) (Catalog, error)
}

type Catalog struct {
	About       About
	Education   Education
	Skills      Skills
	Experiences []Experience
	Projects    []Project
}

type About struct {
	Text string
}

type Education struct {
	Text string
}

type Skills struct {
	Text string
}

type Experience struct {
	Slug string
	Text string
}

type Project struct {
	Slug string
	Text string
}
