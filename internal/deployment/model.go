package deployment

import "time"

type Deployment struct {
	Status     string            `json:"status"`
	Version    string            `json:"version"`
	DeployedAt time.Time         `json:"deployed_at"`
	Services   map[string]string `json:"services"`
}
