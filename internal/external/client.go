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
	atlasURL string
	forgeURL string
	http     *http.Client
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

func NewClient(atlasURL, forgeURL string) *Client {
	return &Client{
		atlasURL: strings.TrimRight(atlasURL, "/"),
		forgeURL: strings.TrimRight(forgeURL, "/"),
		http:     &http.Client{Timeout: 300 * time.Millisecond},
	}
}

func (c *Client) Search(query string) ([]command.SearchResult, error) {
	if c.atlasURL == "" {
		return nil, fmt.Errorf("Atlas URL is not configured")
	}

	endpoint := c.atlasURL + "/search?" + url.Values{"q": {query}}.Encode()
	response, err := c.http.Get(endpoint)
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

	response, err := c.http.Get(c.forgeURL + "/dashboard")
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

	response, err := c.http.Post(
		c.forgeURL+"/events",
		"application/json",
		bytes.NewReader(body),
	)
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
	return c.healthy(c.atlasURL), c.healthy(c.forgeURL)
}

func (c *Client) getJSON(baseURL, path string, target any) error {
	if baseURL == "" {
		return fmt.Errorf("service URL is not configured")
	}
	response, err := c.http.Get(baseURL + path)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("service returned %s", response.Status)
	}
	return json.NewDecoder(response.Body).Decode(target)
}

func (c *Client) healthy(baseURL string) bool {
	if baseURL == "" {
		return false
	}
	response, err := c.http.Get(baseURL + "/healthz")
	if err != nil {
		return false
	}
	defer response.Body.Close()
	return response.StatusCode == http.StatusOK
}
