package command

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
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

type ServiceHealth struct {
	Name   string
	Online bool
	Uptime time.Duration
}

type SystemStatus struct {
	Uptime             time.Duration
	ContentBytes       int64
	ContentBytesKnown  bool
	Services           []ServiceHealth
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

func (Dashboard) Usage() string {
	return "dashboard"
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
		status := d.system.SystemStatus()
		if status.ContentBytesKnown {
			output = addStorageLine(
				output,
				fmt.Sprintf("content          %s", formatBytes(status.ContentBytes)),
			)
		}
		output += "\n\n" + formatSystemStatus(status)
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

func addStorageLine(output string, line string) string {
	const storageHeader = "STORAGE\n=======\n"
	index := strings.Index(output, storageHeader)
	if index == -1 {
		return output + "\n\nSTORAGE\n=======\n" + line
	}

	insertAt := index + len(storageHeader)
	end := strings.Index(output[insertAt:], "\n\n")
	if end == -1 {
		return output + "\n" + line
	}
	return output[:insertAt+end] + "\n" + line + output[insertAt+end:]
}

func formatBytes(bytes int64) string {
	if bytes < 0 {
		bytes = 0
	}
	units := []string{"B", "KiB", "MiB", "GiB", "TiB"}
	value := float64(bytes)
	unit := 0
	for value >= 999.5 && unit < len(units)-1 {
		value /= 1024
		unit++
	}
	if unit == 0 {
		return fmt.Sprintf("%d B", bytes)
	}
	return fmt.Sprintf("%s %s", threeSignificantDigits(value), units[unit])
}

func threeSignificantDigits(value float64) string {
	switch {
	case value >= 100:
		return fmt.Sprintf("%.0f", value)
	case value >= 10:
		return trimTrailingZeros(fmt.Sprintf("%.1f", value))
	default:
		return trimTrailingZeros(fmt.Sprintf("%.2f", math.Max(value, 0.1)))
	}
}

func trimTrailingZeros(value string) string {
	value = strings.TrimRight(value, "0")
	return strings.TrimRight(value, ".")
}

func formatSystemStatus(status SystemStatus) string {
	lines := []string{"SERVICES", "========"}
	for _, service := range status.Services {
		state := "offline"
		detail := ""
		if service.Online {
			state = "online"
			uptime := service.Uptime
			if uptime == 0 && service.Name == "vaultsh" {
				uptime = status.Uptime
			}
			if uptime > 0 {
				detail = "uptime " + formatUptime(uptime)
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
