package sentinel

import (
	"encoding/json"
	"fmt"
	"os"
)

type FileReader struct {
	path string
}

func NewFileReader(path string) *FileReader {
	return &FileReader{path: path}
}

func (r *FileReader) CurrentAssessment() (Assessment, error) {
	var assessment Assessment
	data, err := os.ReadFile(r.path)
	if err != nil {
		return assessment, err
	}
	if err := json.Unmarshal(data, &assessment); err != nil {
		return assessment, fmt.Errorf("decode Sentinel metadata: %w", err)
	}
	if assessment.Commit == "" ||
		assessment.AnalyzedAt.IsZero() ||
		assessment.Risk == "" ||
		assessment.Decision == "" ||
		assessment.Provider == "" ||
		len(assessment.Checks) == 0 {
		return assessment, fmt.Errorf("Sentinel metadata is incomplete")
	}
	return assessment, nil
}
