package deployment

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileReaderReadsDeployment(t *testing.T) {
	path := filepath.Join(t.TempDir(), "deployment.json")
	data := []byte(
		`{"status":"success","version":"deploy-42","deployed_at":"2026-07-03T15:20:00Z"}`,
	)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}

	deployment, err := NewFileReader(path).CurrentDeployment()

	if err != nil {
		t.Fatal(err)
	}
	if deployment.Status != "success" || deployment.Version != "deploy-42" {
		t.Errorf("deployment = %#v", deployment)
	}
}

func TestFileReaderRejectsIncompleteDeployment(t *testing.T) {
	path := filepath.Join(t.TempDir(), "deployment.json")
	if err := os.WriteFile(path, []byte(`{"status":"success"}`), 0o600); err != nil {
		t.Fatal(err)
	}

	if _, err := NewFileReader(path).CurrentDeployment(); err == nil {
		t.Fatal("expected incomplete metadata error")
	}
}
