package sentinel

import (
	"testing"
	"time"
)

func TestFormatAssessment(t *testing.T) {
	assessment := Assessment{
		Commit:     "abcdef123456",
		AnalyzedAt: time.Date(2026, 7, 4, 12, 0, 0, 0, time.UTC),
		Risk:       "high",
		Decision:   "advisory",
		Checks: []AssessmentCheck{
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
	}

	want := "SENTINEL\n========\n" +
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
	if got := FormatAssessment(assessment); got != want {
		t.Errorf("FormatAssessment() = %q", got)
	}
}
