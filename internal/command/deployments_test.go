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

	result := NewDashboard(metrics, deployments).Execute(nil, Input{})

	if result.ExitCode != ExitSuccess ||
		result.Output != "Forge dashboard\n\n"+FormatDeployment(deployments.deployment) {
		t.Errorf("result = %#v", result)
	}
}
