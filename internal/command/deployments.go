package command

import (
	"fmt"
	"strings"

	"github.com/mateom/vaultsh/internal/deployment"
)

type DeploymentService interface {
	CurrentDeployment() (deployment.Deployment, error)
}

type Deployments struct {
	service DeploymentService
}

func NewDeployments(service DeploymentService) Deployments {
	return Deployments{service: service}
}

func (Deployments) Name() string {
	return "deployments"
}

func (Deployments) Description() string {
	return "Show the latest production deployment"
}

func (Deployments) Usage() string {
	return "deployments"
}

func (Deployments) Help() string {
	return "Shows sanitized deployment status, version, and timestamp."
}

func (d Deployments) Execute(args []string, _ Input) Result {
	if len(args) != 0 {
		return Result{Output: "usage: deployments", ExitCode: ExitUsage}
	}
	deployment, err := d.service.CurrentDeployment()
	if err != nil {
		return Result{Output: "deployment status unavailable", ExitCode: ExitFailure}
	}
	return Result{Output: FormatDeployment(deployment), ExitCode: ExitSuccess}
}

func FormatDeployment(deployment deployment.Deployment) string {
	return strings.Join([]string{
		"DEPLOYMENT",
		"==========",
		fmt.Sprintf("  status:  %s", deployment.Status),
		fmt.Sprintf("  version: %s", deployment.Version),
		fmt.Sprintf("  vault:   %s", shortVersion(deployment.Services["vault"])),
		fmt.Sprintf("  atlas:   %s", shortVersion(deployment.Services["atlas"])),
		fmt.Sprintf("  forge:   %s", shortVersion(deployment.Services["forge"])),
		fmt.Sprintf(
			"  updated: %s",
			deployment.DeployedAt.UTC().Format("2006-01-02 15:04:05 UTC"),
		),
	}, "\n")
}

func shortVersion(version string) string {
	if len(version) > 7 {
		return version[:7]
	}
	return version
}
