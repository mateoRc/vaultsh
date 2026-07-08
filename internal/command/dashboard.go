package command

import (
	"fmt"
	"strings"
	"time"
)

type ServiceHealth struct {
	Name   string
	Online bool
	Uptime time.Duration
}

type SystemStatus struct {
	Uptime   time.Duration
	Services []ServiceHealth
}

type SystemService interface {
	SystemStatus() SystemStatus
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
