package command

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type MetricsSummary struct {
	Requests int            `json:"requests"`
	Errors   int            `json:"errors"`
	Average  float64        `json:"avg_ms"`
	Median   float64        `json:"median_ms"`
	Services map[string]int `json:"services"`
	Commands map[string]int `json:"commands"`
}

type MetricsService interface {
	Summary() (MetricsSummary, error)
	Dashboard() (string, error)
}

type ServiceHealth struct {
	Name      string
	Online    bool
	LatencyMS int64
}

type SystemStatus struct {
	Uptime   time.Duration
	Services []ServiceHealth
}

type SystemService interface {
	SystemStatus() SystemStatus
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
		fmt.Sprintf("median ms: %g", summary.Median),
		"services:",
	}
	lines = append(lines, formatCounts(summary.Services)...)
	return Result{Output: strings.Join(lines, "\n"), ExitCode: ExitSuccess}
}

type Dashboard struct {
	metrics    MetricsService
	deployment DeploymentService
	system     SystemService
	assessment AssessmentService
}

func NewDashboard(
	metrics MetricsService,
	deployment DeploymentService,
	system SystemService,
	assessment AssessmentService,
) Dashboard {
	return Dashboard{
		metrics:    metrics,
		deployment: deployment,
		system:     system,
		assessment: assessment,
	}
}

func (Dashboard) Name() string {
	return "dashboard"
}

func (Dashboard) Description() string {
	return "Show the Forge analytics dashboard"
}

func (Dashboard) Help() string {
	return "Renders persisted activity, system status, deployment, and Sentinel data."
}

func (d Dashboard) Execute(args []string, _ Input) Result {
	if len(args) != 0 {
		return Result{Output: "usage: dashboard", ExitCode: ExitUsage}
	}
	output, err := d.metrics.Dashboard()
	if err != nil {
		return Result{Output: "dashboard unavailable", ExitCode: ExitFailure}
	}
	if d.system != nil {
		output += "\n\n" + formatSystemStatus(d.system.SystemStatus())
	}
	if d.deployment != nil {
		current, deploymentErr := d.deployment.CurrentDeployment()
		if deploymentErr == nil {
			output += "\n\n" + FormatDeployment(current)
		}
	}
	if d.assessment != nil {
		current, assessmentErr := d.assessment.CurrentAssessment()
		if assessmentErr == nil {
			output += "\n\n" + FormatAssessment(current)
		}
	}
	return Result{Output: output, ExitCode: ExitSuccess}
}

func formatSystemStatus(status SystemStatus) string {
	lines := []string{"SERVICES", "========"}
	for _, service := range status.Services {
		state := "offline"
		detail := ""
		if service.Online {
			state = "online"
			if service.Name == "vaultsh" {
				detail = "uptime " + formatUptime(status.Uptime)
			} else {
				detail = fmt.Sprintf("latency %d ms", service.LatencyMS)
			}
		}
		line := fmt.Sprintf("  %-8s %s", service.Name, state)
		if detail != "" {
			line += "  " + detail
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func formatUptime(duration time.Duration) string {
	duration = duration.Truncate(time.Second)
	days := duration / (24 * time.Hour)
	duration %= 24 * time.Hour
	if days > 0 {
		return fmt.Sprintf("%dd %s", days, duration)
	}
	return duration.String()
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
