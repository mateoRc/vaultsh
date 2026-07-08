package command

import (
	"fmt"
	"sort"
	"strings"
)

type MetricsSummary struct {
	Requests   int            `json:"requests"`
	Errors     int            `json:"errors"`
	UserErrors int            `json:"user_errors"`
	Average    float64        `json:"avg_ms"`
	Median     float64        `json:"median_ms"`
	Services   map[string]int `json:"services"`
	Commands   map[string]int `json:"commands"`
}

type MetricsService interface {
	Summary() (MetricsSummary, error)
	Dashboard() (string, error)
}

type Metrics struct {
	service MetricsService
}

func NewMetrics(service MetricsService) Metrics {
	return Metrics{service: service}
}

func (Metrics) Name() string {
	return "metrics"
}

func (Metrics) Description() string {
	return "Show service analytics from Forge"
}

func (Metrics) Usage() string {
	return "metrics"
}

func (Metrics) Help() string {
	return "Shows request, error, response time, service, and command totals."
}

func (m Metrics) Execute(args []string, _ Input) Result {
	if len(args) != 0 {
		return Result{Output: "usage: metrics", ExitCode: ExitUsage}
	}
	summary, err := m.service.Summary()
	if err != nil {
		return Result{Output: "metrics unavailable", ExitCode: ExitFailure}
	}

	lines := []string{
		fmt.Sprintf("requests: %d", summary.Requests),
		fmt.Sprintf("runtime errors: %d", summary.Errors),
		fmt.Sprintf("user errors: %d", summary.UserErrors),
		fmt.Sprintf("avg ms: %g", summary.Average),
		fmt.Sprintf("median ms: %g", summary.Median),
		"services:",
	}
	lines = append(lines, formatCounts(summary.Services)...)
	return Result{Output: strings.Join(lines, "\n"), ExitCode: ExitSuccess}
}

func formatCounts(counts map[string]int) []string {
	if len(counts) == 0 {
		return []string{"  none"}
	}
	names := make([]string, 0, len(counts))
	for name := range counts {
		names = append(names, name)
	}
	sort.Strings(names)
	lines := make([]string, len(names))
	for index, name := range names {
		lines[index] = fmt.Sprintf("  %s: %d", name, counts[name])
	}
	return lines
}
