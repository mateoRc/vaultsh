package command

import (
	"fmt"
	"strings"
	"time"

	"github.com/mateom/vaultsh/internal/sentinel"
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

type AssessmentService interface {
	CurrentAssessment() (sentinel.Assessment, error)
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
	return "dashboard [--window 24h|7d|30d|all] [--include-invalid]"
}

func (Dashboard) Help() string {
	return "Defaults to the last 24 hours; all includes every retained event. Unknown command names are hidden unless --include-invalid is set."
}

func (d Dashboard) Execute(args []string, _ Input) Result {
	windowHours, includeInvalid, ok := dashboardOptions(args)
	if !ok {
		return Result{Output: "usage: " + d.Usage(), ExitCode: ExitUsage}
	}
	output, err := d.metrics.Dashboard(windowHours, includeInvalid)
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
			output += "\n\n" + sentinel.FormatAssessment(current)
		}
	}
	return Result{Output: output, ExitCode: ExitSuccess}
}

func dashboardOptions(args []string) (int, bool, bool) {
	windowHours := 24
	includeInvalid := false
	windowSet := false
	for index := 0; index < len(args); index++ {
		switch args[index] {
		case "--include-invalid":
			if includeInvalid {
				return 0, false, false
			}
			includeInvalid = true
		case "--window":
			if windowSet || index+1 >= len(args) {
				return 0, false, false
			}
			index++
			windowSet = true
			switch args[index] {
			case "24h":
				windowHours = 24
			case "7d":
				windowHours = 168
			case "30d":
				windowHours = 720
			case "all":
				windowHours = 0
			default:
				return 0, false, false
			}
		default:
			return 0, false, false
		}
	}
	return windowHours, includeInvalid, true
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
