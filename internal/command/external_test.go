package command

import (
	"errors"
	"testing"
)

type externalStub struct {
	searchResults []SearchResult
	searchErr     error
	invokeAnswer  Answer
	invokeErr     error
	askAnswer     Answer
	askErr        error
	summary       MetricsSummary
	dashboard     string
	metricsErr    error
	dashboardArgs *struct {
		windowHours    int
		includeInvalid bool
	}
}

func (s externalStub) Search(string) ([]SearchResult, error) {
	return s.searchResults, s.searchErr
}

func (s externalStub) Summary() (MetricsSummary, error) {
	return s.summary, s.metricsErr
}

func (s externalStub) Dashboard(windowHours int, includeInvalid bool) (string, error) {
	if s.dashboardArgs != nil {
		s.dashboardArgs.windowHours = windowHours
		s.dashboardArgs.includeInvalid = includeInvalid
	}
	return s.dashboard, s.metricsErr
}

func (s externalStub) Invoke(string) (Answer, error) {
	return s.invokeAnswer, s.invokeErr
}

func (s externalStub) Ask(string) (Answer, error) {
	return s.askAnswer, s.askErr
}

func TestDashboardWindowAndInvalidFlag(t *testing.T) {
	for _, test := range []struct {
		name           string
		args           []string
		windowHours    int
		includeInvalid bool
		exitCode       int
	}{
		{name: "default", windowHours: 24},
		{name: "week", args: []string{"--window", "7d"}, windowHours: 168},
		{name: "month", args: []string{"--window", "30d"}, windowHours: 720},
		{name: "all retained", args: []string{"--window", "all"}, windowHours: 0},
		{name: "include invalid", args: []string{"--include-invalid", "--window", "all"}, includeInvalid: true},
		{name: "bad window", args: []string{"--window", "forever"}, exitCode: ExitUsage},
		{name: "duplicate flag", args: []string{"--include-invalid", "--include-invalid"}, exitCode: ExitUsage},
	} {
		t.Run(test.name, func(t *testing.T) {
			seen := &struct {
				windowHours    int
				includeInvalid bool
			}{}
			result := NewDashboard(externalStub{dashboard: "Forge dashboard", dashboardArgs: seen}, nil, nil, nil).
				Execute(test.args, Input{})
			if result.ExitCode != test.exitCode {
				t.Fatalf("exit code = %d, want %d", result.ExitCode, test.exitCode)
			}
			if test.exitCode == ExitSuccess &&
				(seen.windowHours != test.windowHours || seen.includeInvalid != test.includeInvalid) {
				t.Fatalf("dashboard options = %+v", seen)
			}
		})
	}
}

func TestSearchFormatsAtlasResults(t *testing.T) {
	search := NewSearch(externalStub{searchResults: []SearchResult{{
		Path:       "/cv/skills.md",
		LineNumber: 12,
		Line:       "messaging: Kafka",
	}}})

	result := search.Execute([]string{"kafka"}, Input{})

	if result.ExitCode != ExitSuccess {
		t.Fatalf("exit code = %d, want %d", result.ExitCode, ExitSuccess)
	}
	if result.Output != "/cv/skills.md:12: messaging: Kafka" {
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

func TestInvokeFormatsAnswerAndSources(t *testing.T) {
	invoke := NewInvoke(externalStub{invokeAnswer: Answer{
		Answer:  "Vaultsh is the Go shell.",
		Sources: []string{"/docs/projects/overview.md", "/cv/about.md"},
	}})

	result := invoke.Execute([]string{"what", "is", "vaultsh"}, Input{})

	if result.ExitCode != ExitSuccess {
		t.Fatalf("exit code = %d, want %d", result.ExitCode, ExitSuccess)
	}
	want := "Vaultsh is the Go shell.\n\nsources: /docs/projects/overview.md, /cv/about.md"
	if result.Output != want {
		t.Errorf("output = %q, want %q", result.Output, want)
	}
}

func TestInvokeRequiresAQuestion(t *testing.T) {
	result := NewInvoke(externalStub{}).Execute(nil, Input{})

	if result.ExitCode != ExitUsage || result.Output != "usage: invoke <question>" {
		t.Errorf("result = %#v", result)
	}
}

func TestInvokeDegradesWhenOracleIsUnavailable(t *testing.T) {
	invoke := NewInvoke(externalStub{invokeErr: errors.New("offline")})

	result := invoke.Execute([]string{"who", "is", "mateo"}, Input{})

	if result.ExitCode != ExitFailure || result.Output != oracleWithdrawn {
		t.Errorf("result = %#v", result)
	}
}

func TestAskFormatsAnswerAndSources(t *testing.T) {
	ask := NewAsk(externalStub{askAnswer: Answer{
		Answer:  "Mateo writes Go and Python.",
		Sources: []string{"/cv/skills.md"},
	}})

	result := ask.Execute([]string{"which", "languages"}, Input{})

	if result.ExitCode != ExitSuccess {
		t.Fatalf("exit code = %d, want %d", result.ExitCode, ExitSuccess)
	}
	want := "Mateo writes Go and Python.\n\nsources: /cv/skills.md"
	if result.Output != want {
		t.Errorf("output = %q, want %q", result.Output, want)
	}
}

func TestAskRequiresAQuestion(t *testing.T) {
	result := NewAsk(externalStub{}).Execute(nil, Input{})

	if result.ExitCode != ExitUsage || result.Output != "usage: ask <question>" {
		t.Errorf("result = %#v", result)
	}
}

func TestAskDegradesWhenOracleIsUnavailable(t *testing.T) {
	ask := NewAsk(externalStub{askErr: errors.New("offline")})

	result := ask.Execute([]string{"hello"}, Input{})

	if result.ExitCode != ExitFailure || result.Output != assistantOffline {
		t.Errorf("result = %#v", result)
	}
}

func TestMetricsAndDashboardFormatForgeResponses(t *testing.T) {
	service := externalStub{
		summary: MetricsSummary{
			Requests:   3,
			Errors:     1,
			UserErrors: 1,
			Average:    6,
			Median:     4,
			Services:   map[string]int{"vault": 2, "atlas": 1},
		},
		dashboard: "Forge dashboard",
	}

	metricsResult := NewMetrics(service).Execute(nil, Input{})
	if metricsResult.ExitCode != ExitSuccess {
		t.Fatalf("metrics exit code = %d", metricsResult.ExitCode)
	}
	if metricsResult.Output != "requests: 3\nruntime errors: 1\nuser errors: 1\navg ms: 6\nmedian ms: 4\nservices:\n  atlas: 1\n  vault: 2" {
		t.Errorf("metrics output = %q", metricsResult.Output)
	}

	dashboardResult := NewDashboard(service, nil, nil, nil).Execute(nil, Input{})
	if dashboardResult.Output != "Forge dashboard" {
		t.Errorf("dashboard output = %q", dashboardResult.Output)
	}
}
