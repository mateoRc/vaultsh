package deployment

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

func (r *FileReader) CurrentDeployment() (Deployment, error) {
	var deployment Deployment
	data, err := os.ReadFile(r.path)
	if err != nil {
		return deployment, err
	}
	if err := json.Unmarshal(data, &deployment); err != nil {
		return deployment, fmt.Errorf("decode deployment metadata: %w", err)
	}
	if deployment.Status == "" ||
		deployment.Version == "" ||
		deployment.DeployedAt.IsZero() ||
		deployment.Services["vault"] == "" ||
		deployment.Services["atlas"] == "" ||
		deployment.Services["forge"] == "" {
		return deployment, fmt.Errorf("deployment metadata is incomplete")
	}
	return deployment, nil
}
