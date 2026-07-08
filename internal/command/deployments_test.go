package command

import (
	"errors"
	"testing"
	"time"

	"github.com/mateom/vaultsh/internal/deployment"
	"github.com/mateom/vaultsh/internal/sentinel"
)

type deploymentStub struct {
	deployment deployment.Deployment
	err        error
}

type systemStub struct {
	status SystemStatus
}

type assessmentStub struct {
	assessment sentinel.Assessment
	err        error
}

func (s systemStub) SystemStatus() SystemStatus {
	return s.status
}

func (s assessmentStub) CurrentAssessment() (sentinel.Assessment, error) {
	return s.assessment, s.err
}

func (s deploymentStub) CurrentDeployment() (deployment.Deployment, error) {
	return s.deployment, s.err
}

func TestDeploymentsFormatsSanitizedMetadata(t *testing.T) {
	service := deploymentStub{deployment: deployment.Deployment{
		Status:     "success",
		Version:    "deploy-42",
		DeployedAt: time.Date(2026, 7, 3, 15, 20, 0, 0, time.UTC),
		Services: map[string]string{
			"vault": "111111111111",
			"atlas": "222222222222",
			"forge": "333333333333",
		},
	}}

	result := NewDeployments(service).Execute(nil, Input{})

	want := "DEPLOYMENT\n" +
		"==========\n" +
		"  status:  success\n" +
		"  version: deploy-42\n" +
		"  vault:   1111111\n" +
		"  atlas:   2222222\n" +
		"  forge:   3333333\n" +
		"  updated: 2026-07-03 15:20:00 UTC"
	if result.ExitCode != ExitSuccess || result.Output != want {
		t.Errorf("result = %#v", result)
	}
}

func TestDeploymentsDegradesWhenMetadataIsUnavailable(t *testing.T) {
	result := NewDeployments(deploymentStub{err: errors.New("missing")}).
		Execute(nil, Input{})

	if result.ExitCode != ExitFailure || result.Output != "deployment status unavailable" {
		t.Errorf("result = %#v", result)
	}
}

func TestDashboardIncludesDeployment(t *testing.T) {
	metrics := externalStub{dashboard: "Forge dashboard"}
	deployments := deploymentStub{deployment: deployment.Deployment{
		Status:     "success",
		Version:    "deploy-42",
		DeployedAt: time.Date(2026, 7, 3, 15, 20, 0, 0, time.UTC),
		Services: map[string]string{
			"vault": "1111111",
			"atlas": "2222222",
			"forge": "3333333",
		},
	}}

	result := NewDashboard(metrics, deployments, nil, nil).Execute(nil, Input{})

	if result.ExitCode != ExitSuccess ||
		result.Output != "Forge dashboard\n\n"+FormatDeployment(deployments.deployment) {
		t.Errorf("result = %#v", result)
	}
}

func TestDashboardIncludesServiceHealthAndUptime(t *testing.T) {
	metrics := externalStub{dashboard: "Forge dashboard"}
	system := systemStub{status: SystemStatus{
		Uptime: 26*time.Hour + 3*time.Minute,
		Services: []ServiceHealth{
			{Name: "vaultsh", Online: true},
			{Name: "atlas", Online: true, Uptime: time.Hour},
			{Name: "forge", Online: true, Uptime: 2 * time.Hour},
		},
	}}

	result := NewDashboard(metrics, nil, system, nil).Execute(nil, Input{})

	want := "Forge dashboard\n\nSERVICES\n========\n" +
		"  vaultsh  online  uptime 1d 2h3m0s\n" +
		"  atlas    online  uptime 1h0m0s\n" +
		"  forge    online  uptime 2h0m0s"
	if result.ExitCode != ExitSuccess || result.Output != want {
		t.Errorf("result = %#v", result)
	}
}

func TestDashboardIncludesSentinelAssessment(t *testing.T) {
	metrics := externalStub{dashboard: "Forge dashboard"}
	assessment := assessmentStub{assessment: sentinel.Assessment{
		Commit:     "abcdef123456",
		AnalyzedAt: time.Date(2026, 7, 4, 12, 0, 0, 0, time.UTC),
		Risk:       "high",
		Decision:   "advisory",
		Checks: []sentinel.AssessmentCheck{
			{Name: "tests", Status: "passed"},
			{
				Name:     "vaultsh-image-security",
				Status:   "failed",
				Source:   "Trivy",
				Evidence: "1 high vulnerability; CVE-1 lib 1.0 -> 1.1",
			},
		},
		Summary: "1 of 2 checks failed: vaultsh image security.",
		Actions: []string{
			"Update affected packages to listed fixed versions.",
		},
		Provider: "mock",
	}}

	result := NewDashboard(metrics, nil, nil, assessment).Execute(nil, Input{})

	want := "Forge dashboard\n\nSENTINEL\n========\n" +
		"  decision: advisory\n" +
		"  risk:     high\n" +
		"  checks:   1 passed, 0 warning, 1 failed\n" +
		"  provider: mock\n" +
		"  commit:   abcdef1\n" +
		"  updated:  2026-07-04 12:00:00 UTC\n" +
		"  summary:  1 of 2 checks failed: vaultsh image security.\n" +
		"  findings:\n" +
		"    - [failed] vaultsh image security (Trivy): " +
		"1 high vulnerability; CVE-1 lib 1.0 -> 1.1\n" +
		"  next actions:\n" +
		"    - Update affected packages to listed fixed versions."
	if result.ExitCode != ExitSuccess || result.Output != want {
		t.Errorf("result = %#v", result)
	}
}
