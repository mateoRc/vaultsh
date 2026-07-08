package sentinel

import (
	"fmt"
	"strings"
)

func FormatAssessment(assessment Assessment) string {
	counts := countCheckStatuses(assessment.Checks)
	lines := []string{
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
		fmt.Sprintf("  commit:   %s", shortCommit(assessment.Commit)),
		fmt.Sprintf(
			"  updated:  %s",
			assessment.AnalyzedAt.UTC().Format("2006-01-02 15:04:05 UTC"),
		),
		fmt.Sprintf("  summary:  %s", assessment.Summary),
	}
	lines = append(lines, formatFindings(assessment.Checks)...)
	lines = append(lines, formatActions(assessment.Actions)...)
	return strings.Join(lines, "\n")
}

func countCheckStatuses(checks []AssessmentCheck) map[string]int {
	counts := map[string]int{}
	for _, check := range checks {
		counts[check.Status]++
	}
	return counts
}

func shortCommit(commit string) string {
	if len(commit) <= 7 {
		return commit
	}
	return commit[:7]
}

func formatFindings(checks []AssessmentCheck) []string {
	lines := []string{}
	for _, check := range checks {
		if check.Status == "passed" {
			continue
		}
		if len(lines) == 0 {
			lines = append(lines, "  findings:")
		}
		lines = append(lines, formatFinding(check))
	}
	return lines
}

func formatFinding(check AssessmentCheck) string {
	return fmt.Sprintf(
		"    - [%s] %s (%s): %s",
		check.Status,
		strings.ReplaceAll(check.Name, "-", " "),
		check.Source,
		check.Evidence,
	)
}

func formatActions(actions []string) []string {
	if len(actions) == 0 {
		return nil
	}

	lines := []string{"  next actions:"}
	for _, action := range actions {
		lines = append(lines, "    - "+action)
	}
	return lines
}
