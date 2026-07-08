package sentinel

import "time"

type AssessmentCheck struct {
	Name     string `json:"name"`
	Status   string `json:"status"`
	Source   string `json:"source"`
	Evidence string `json:"evidence"`
}

type Assessment struct {
	Commit     string            `json:"commit"`
	AnalyzedAt time.Time         `json:"analyzed_at"`
	Risk       string            `json:"risk"`
	Decision   string            `json:"decision"`
	Checks     []AssessmentCheck `json:"checks"`
	Summary    string            `json:"summary"`
	Actions    []string          `json:"actions"`
	Provider   string            `json:"provider"`
}
