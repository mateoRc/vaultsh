package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/mateom/vaultsh/internal/command"
)

type Client struct {
	atlasURL   string
	forgeURL   string
	atlasToken string
	forgeToken string
	http       *http.Client
	startedAt  time.Time
}

type searchResponse struct {
	Results []command.SearchResult `json:"results"`
}

type event struct {
	Service    string `json:"service"`
	Event      string `json:"event"`
	Name       string `json:"name"`
	DurationMS int64  `json:"duration_ms"`
	ExitCode   int    `json:"exit_code"`
}

type serviceStatusResponse struct {
	UptimeSeconds int64 `json:"uptime_seconds"`
}

func NewClient(atlasURL, forgeURL, atlasToken, forgeToken string) *Client {
	return &Client{
		atlasURL:   strings.TrimRight(atlasURL, "/"),
		forgeURL:   strings.TrimRight(forgeURL, "/"),
		atlasToken: atlasToken,
		forgeToken: forgeToken,
		http:       &http.Client{Timeout: 300 * time.Millisecond},
		startedAt:  time.Now(),
	}
}

func (c *Client) Search(query string) ([]command.SearchResult, error) {
	if c.atlasURL == "" {
		return nil, fmt.Errorf("Atlas URL is not configured")
	}

	endpoint := c.atlasURL + "/search?" + url.Values{"q": {query}}.Encode()
	response, err := c.get(endpoint, c.atlasToken)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Atlas returned %s", response.Status)
	}

	var result searchResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Results, nil
}

func (c *Client) Summary() (command.MetricsSummary, error) {
	var summary command.MetricsSummary
	err := c.getJSON(c.forgeURL, "/summary", &summary)
	return summary, err
}

func (c *Client) Dashboard() (string, error) {
	if c.forgeURL == "" {
		return "", fmt.Errorf("Forge URL is not configured")
	}

	response, err := c.get(c.forgeURL+"/dashboard", c.forgeToken)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Forge returned %s", response.Status)
	}

	var output bytes.Buffer
	if _, err := output.ReadFrom(response.Body); err != nil {
		return "", err
	}
	return output.String(), nil
}

func (c *Client) Record(
	service string,
	eventName string,
	name string,
	durationMS int64,
	exitCode int,
) error {
	if c.forgeURL == "" {
		return nil
	}

	body, err := json.Marshal(event{
		Service:    service,
		Event:      eventName,
		Name:       name,
		DurationMS: durationMS,
		ExitCode:   exitCode,
	})
	if err != nil {
		return err
	}

	request, err := http.NewRequest(
		http.MethodPost,
		c.forgeURL+"/events",
		bytes.NewReader(body),
	)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+c.forgeToken)
	response, err := c.http.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Forge returned %s", response.Status)
	}
	return nil
}

func (c *Client) Availability() (bool, bool) {
	atlas, _ := c.probe(c.atlasURL)
	forge, _ := c.probe(c.forgeURL)
	return atlas, forge
}

func (c *Client) SystemStatus() command.SystemStatus {
	vaultUptime := time.Since(c.startedAt)
	atlasOnline, atlasUptime := c.serviceStatus(c.atlasURL)
	forgeOnline, forgeUptime := c.serviceStatus(c.forgeURL)
	return command.SystemStatus{
		Uptime: vaultUptime,
		Services: []command.ServiceHealth{
			{Name: "vaultsh", Online: true, Uptime: vaultUptime},
			{Name: "atlas", Online: atlasOnline, Uptime: atlasUptime},
			{Name: "forge", Online: forgeOnline, Uptime: forgeUptime},
		},
	}
}

func (c *Client) getJSON(baseURL, path string, target any) error {
	if baseURL == "" {
		return fmt.Errorf("service URL is not configured")
	}
	response, err := c.get(baseURL+path, c.forgeToken)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("service returned %s", response.Status)
	}
	return json.NewDecoder(response.Body).Decode(target)
}

func (c *Client) get(endpoint, token string) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+token)
	return c.http.Do(request)
}

func (c *Client) probe(baseURL string) (bool, time.Duration) {
	if baseURL == "" {
		return false, 0
	}
	startedAt := time.Now()
	response, err := c.http.Get(baseURL + "/healthz")
	if err != nil {
		return false, time.Since(startedAt)
	}
	defer response.Body.Close()
	return response.StatusCode == http.StatusOK, time.Since(startedAt)
}

func (c *Client) serviceStatus(baseURL string) (bool, time.Duration) {
	if baseURL == "" {
		return false, 0
	}

	response, err := c.http.Get(baseURL + "/status")
	if err != nil {
		return false, 0
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		online, _ := c.probe(baseURL)
		return online, 0
	}

	var status serviceStatusResponse
	if err := json.NewDecoder(response.Body).Decode(&status); err != nil {
		online, _ := c.probe(baseURL)
		return online, 0
	}
	return true, time.Duration(status.UptimeSeconds) * time.Second
}
