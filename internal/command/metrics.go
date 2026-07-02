package command

import (
	"fmt"
	"sort"
	"strings"
)

type MetricsSummary struct {
	Requests int            `json:"requests"`
	Errors   int            `json:"errors"`
	Average  float64        `json:"avg_ms"`
	Services map[string]int `json:"services"`
	Commands map[string]int `json:"commands"`
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

func (Metrics) Help() string {
	return "Shows request, error, latency, service, and command totals."
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
		fmt.Sprintf("errors: %d", summary.Errors),
		fmt.Sprintf("avg ms: %g", summary.Average),
		"services:",
	}
	lines = append(lines, formatCounts(summary.Services)...)
	return Result{Output: strings.Join(lines, "\n"), ExitCode: ExitSuccess}
}

type Dashboard struct {
	service MetricsService
}

func NewDashboard(service MetricsService) Dashboard {
	return Dashboard{service: service}
}

func (Dashboard) Name() string {
	return "dashboard"
}

func (Dashboard) Description() string {
	return "Show the Forge analytics dashboard"
}

func (Dashboard) Help() string {
	return "Renders live service activity as a terminal-friendly dashboard."
}

func (d Dashboard) Execute(args []string, _ Input) Result {
	if len(args) != 0 {
		return Result{Output: "usage: dashboard", ExitCode: ExitUsage}
	}
	output, err := d.service.Dashboard()
	if err != nil {
		return Result{Output: "dashboard unavailable", ExitCode: ExitFailure}
	}
	return Result{Output: output, ExitCode: ExitSuccess}
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
