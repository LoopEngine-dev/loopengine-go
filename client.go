package loopengine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const defaultBaseURL = "https://api.loopengine.dev"
const apiPath = "/feedback"

// Client sends feedback to the LoopEngine Ingest API. Safe for concurrent use.
type Client struct {
	projectKey    string
	projectSecret string
	projectID     string
	baseURL       string
	httpClient    *http.Client
}

// New builds a Client from project credentials. Options can customize the HTTP client (e.g. timeouts).
func New(projectKey, projectSecret, projectID string, opts ...Option) (*Client, error) {
	if projectKey == "" || projectSecret == "" || projectID == "" {
		return nil, fmt.Errorf("loopengine: project_key, project_secret, and project_id are required")
	}
	c := &Client{
		projectKey:    strings.TrimSpace(projectKey),
		projectSecret: strings.TrimSpace(projectSecret),
		projectID:     strings.TrimSpace(projectID),
		baseURL:       defaultBaseURL,
		httpClient:    http.DefaultClient,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c, nil
}

// Send posts the payload to the Ingest API. payload is JSON-encoded; project_id is set automatically.
// payload can be a map, struct, or any type that encoding/json can marshal.
func (c *Client) Send(ctx context.Context, payload any) error {
	body, err := c.buildBody(payload)
	if err != nil {
		return err
	}
	timestamp, signature := signRequest(c.projectSecret, http.MethodPost, apiPath, body)
	url := c.baseURL + apiPath
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("loopengine: request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Project-Key", c.projectKey)
	req.Header.Set("X-Timestamp", timestamp)
	req.Header.Set("X-Signature", signature)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("loopengine: send: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("loopengine: %s %s", resp.Status, readBody(resp))
	}
	return nil
}

// buildBody marshals payload to JSON and ensures project_id is set.
func (c *Client) buildBody(payload any) ([]byte, error) {
	var m map[string]any
	switch v := payload.(type) {
	case map[string]any:
		m = make(map[string]any, len(v)+1)
		for k, val := range v {
			m[k] = val
		}
	case nil:
		m = make(map[string]any, 1)
	default:
		b, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("loopengine: marshal payload: %w", err)
		}
		if err := json.Unmarshal(b, &m); err != nil {
			return nil, fmt.Errorf("loopengine: payload must be JSON-serializable: %w", err)
		}
	}
	if m == nil {
		m = make(map[string]any, 1)
	}
	m["project_id"] = c.projectID
	return json.Marshal(m)
}

func readBody(resp *http.Response) string {
	b, _ := io.ReadAll(resp.Body)
	return strings.TrimSpace(string(b))
}
