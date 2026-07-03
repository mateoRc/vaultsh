package command

import (
	"errors"
	"testing"
)

type externalStub struct {
	searchResults []SearchResult
	searchErr     error
	summary       MetricsSummary
	dashboard     string
	metricsErr    error
}

func (s externalStub) Search(string) ([]SearchResult, error) {
	return s.searchResults, s.searchErr
}

func (s externalStub) Summary() (MetricsSummary, error) {
	return s.summary, s.metricsErr
}

func (s externalStub) Dashboard() (string, error) {
	return s.dashboard, s.metricsErr
}

func TestSearchFormatsAtlasResults(t *testing.T) {
	search := NewSearch(externalStub{searchResults: []SearchResult{{
		Path:       "/cv/skills.txt",
		LineNumber: 12,
		Line:       "messaging: Kafka",
	}}})

	result := search.Execute([]string{"kafka"}, Input{})

	if result.ExitCode != ExitSuccess {
		t.Fatalf("exit code = %d, want %d", result.ExitCode, ExitSuccess)
	}
	if result.Output != "/cv/skills.txt:12: messaging: Kafka" {
		t.Errorf("output = %q", result.Output)
	}
}

func TestSearchDegradesWhenAtlasIsUnavailable(t *testing.T) {
	search := NewSearch(externalStub{searchErr: errors.New("offline")})

	result := search.Execute([]string{"kafka"}, Input{})

	if result.ExitCode != ExitFailure || result.Output != "search unavailable" {
		t.Errorf("result = %#v", result)
	}
}

func TestMetricsAndDashboardFormatForgeResponses(t *testing.T) {
	service := externalStub{
		summary: MetricsSummary{
			Requests: 3,
			Errors:   1,
			Average:  6,
			Services: map[string]int{"vault": 2, "atlas": 1},
		},
		dashboard: "Forge dashboard",
	}

	metricsResult := NewMetrics(service).Execute(nil, Input{})
	if metricsResult.ExitCode != ExitSuccess {
		t.Fatalf("metrics exit code = %d", metricsResult.ExitCode)
	}
	if metricsResult.Output != "requests: 3\nerrors: 1\navg ms: 6\nservices:\n  atlas: 1\n  vault: 2" {
		t.Errorf("metrics output = %q", metricsResult.Output)
	}

	dashboardResult := NewDashboard(service, nil, nil).Execute(nil, Input{})
	if dashboardResult.Output != "Forge dashboard" {
		t.Errorf("dashboard output = %q", dashboardResult.Output)
	}
}
