package command

import (
	"errors"
	"testing"
	"time"
)

type deploymentStub struct {
	deployment Deployment
	err        error
}

type systemStub struct {
	status SystemStatus
}

type assessmentStub struct {
	assessment Assessment
	err        error
}

func (s systemStub) SystemStatus() SystemStatus {
	return s.status
}

func (s assessmentStub) CurrentAssessment() (Assessment, error) {
	return s.assessment, s.err
}

func (s deploymentStub) CurrentDeployment() (Deployment, error) {
	return s.deployment, s.err
}

func TestDeploymentsFormatsSanitizedMetadata(t *testing.T) {
	service := deploymentStub{deployment: Deployment{
		Status:     "success",
		Version:    "deploy-42",
		DeployedAt: time.Date(2026, 7, 3, 15, 20, 0, 0, time.UTC),
	}}

	result := NewDeployments(service).Execute(nil, Input{})

	want := "DEPLOYMENT\n" +
		"==========\n" +
		"  status:  success\n" +
		"  version: deploy-42\n" +
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
	deployments := deploymentStub{deployment: Deployment{
		Status:     "success",
		Version:    "deploy-42",
		DeployedAt: time.Date(2026, 7, 3, 15, 20, 0, 0, time.UTC),
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
			{Name: "atlas", Online: true, LatencyMS: 12},
			{Name: "forge", Online: false},
		},
	}}

	result := NewDashboard(metrics, nil, system, nil).Execute(nil, Input{})

	want := "Forge dashboard\n\nSERVICES\n========\n" +
		"  vaultsh  online  uptime 1d 2h3m0s\n" +
		"  atlas    online  latency 12 ms\n" +
		"  forge    offline"
	if result.ExitCode != ExitSuccess || result.Output != want {
		t.Errorf("result = %#v", result)
	}
}

func TestDashboardIncludesSentinelAssessment(t *testing.T) {
	metrics := externalStub{dashboard: "Forge dashboard"}
	assessment := assessmentStub{assessment: Assessment{
		Commit:     "abcdef123456",
		AnalyzedAt: time.Date(2026, 7, 4, 12, 0, 0, 0, time.UTC),
		Risk:       "low",
		Checks: []AssessmentCheck{
			{Name: "tests", Status: "passed"},
			{Name: "security", Status: "warning"},
		},
		Summary:  "Backend Lab passed with warnings.",
		Provider: "mock",
	}}

	result := NewDashboard(metrics, nil, nil, assessment).Execute(nil, Input{})

	want := "Forge dashboard\n\nSENTINEL\n========\n" +
		"  mode:     advisory\n" +
		"  risk:     low\n" +
		"  checks:   1 passed, 1 warning, 0 failed\n" +
		"  provider: mock\n" +
		"  commit:   abcdef1\n" +
		"  updated:  2026-07-04 12:00:00 UTC\n" +
		"  summary:  Backend Lab passed with warnings."
	if result.ExitCode != ExitSuccess || result.Output != want {
		t.Errorf("result = %#v", result)
	}
}
