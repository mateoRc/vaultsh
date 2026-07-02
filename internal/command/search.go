package command

import (
	"fmt"
	"strings"
)

type SearchResult struct {
	Path       string `json:"path"`
	LineNumber int    `json:"line_number"`
	Line       string `json:"line"`
}

type SearchService interface {
	Search(query string) ([]SearchResult, error)
}

type Search struct {
	service SearchService
}

func NewSearch(service SearchService) Search {
	return Search{service: service}
}

func (Search) Name() string {
	return "search"
}

func (Search) Description() string {
	return "Search portfolio content with the Atlas search engine"
}

func (Search) Usage() string {
	return "search <query>"
}

func (Search) Help() string {
	return "Example: search distributed systems"
}

func (s Search) Execute(args []string, _ Input) Result {
	if len(args) == 0 {
		return Result{Output: "usage: search <query>", ExitCode: ExitUsage}
	}

	results, err := s.service.Search(strings.Join(args, " "))
	if err != nil {
		return Result{Output: "search unavailable", ExitCode: ExitFailure}
	}
	if len(results) == 0 {
		return Result{Output: "no matches", ExitCode: ExitSuccess}
	}

	lines := make([]string, len(results))
	for index, result := range results {
		lines[index] = fmt.Sprintf(
			"%s:%d: %s",
			result.Path,
			result.LineNumber,
			result.Line,
		)
	}
	return Result{Output: strings.Join(lines, "\n"), ExitCode: ExitSuccess}
}
