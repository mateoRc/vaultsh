package sentinel

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileReaderLoadsAssessment(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sentinel.json")
	data := `{
		"commit":"abcdef123456",
		"analyzed_at":"2026-07-04T12:00:00Z",
		"risk":"low",
		"checks":[{"name":"tests","status":"passed"}],
		"summary":"All checks passed.",
		"provider":"mock"
	}`
	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatal(err)
	}

	assessment, err := NewFileReader(path).CurrentAssessment()

	if err != nil {
		t.Fatal(err)
	}
	if assessment.Risk != "low" || assessment.Provider != "mock" {
		t.Errorf("assessment = %#v", assessment)
	}
}

func TestFileReaderRejectsIncompleteAssessment(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sentinel.json")
	if err := os.WriteFile(path, []byte(`{"risk":"low"}`), 0o600); err != nil {
		t.Fatal(err)
	}

	if _, err := NewFileReader(path).CurrentAssessment(); err == nil {
		t.Fatal("expected incomplete metadata error")
	}
}
