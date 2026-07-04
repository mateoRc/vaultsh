package command

import (
	"fmt"
	"strings"
	"time"
)

type AssessmentCheck struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type Assessment struct {
	Commit     string            `json:"commit"`
	AnalyzedAt time.Time         `json:"analyzed_at"`
	Risk       string            `json:"risk"`
	Decision   string            `json:"decision"`
	Checks     []AssessmentCheck `json:"checks"`
	Summary    string            `json:"summary"`
	Provider   string            `json:"provider"`
}

type AssessmentService interface {
	CurrentAssessment() (Assessment, error)
}

func FormatAssessment(assessment Assessment) string {
	counts := map[string]int{}
	for _, check := range assessment.Checks {
		counts[check.Status]++
	}
	commit := assessment.Commit
	if len(commit) > 7 {
		commit = commit[:7]
	}
	return strings.Join([]string{
		"SENTINEL",
		"========",
		fmt.Sprintf("  decision: %s", assessment.Decision),
		fmt.Sprintf("  risk:     %s", assessment.Risk),
		fmt.Sprintf(
			"  checks:   %d passed, %d warning, %d failed",
			counts["passed"],
			counts["warning"],
			counts["failed"],
		),
		fmt.Sprintf("  provider: %s", assessment.Provider),
		fmt.Sprintf("  commit:   %s", commit),
		fmt.Sprintf(
			"  updated:  %s",
			assessment.AnalyzedAt.UTC().Format("2006-01-02 15:04:05 UTC"),
		),
		fmt.Sprintf("  summary:  %s", assessment.Summary),
	}, "\n")
}
