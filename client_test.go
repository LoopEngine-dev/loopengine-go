package loopengine

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// redirectTransport sends requests to target (used in tests only).
type redirectTransport struct {
	target *url.URL
	rt     http.RoundTripper
}

func (r redirectTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	req.URL.Scheme = r.target.Scheme
	req.URL.Host = r.target.Host
	return r.rt.RoundTrip(req)
}

func TestNew_validation(t *testing.T) {
	tests := []struct {
		name   string
		key    string
		secret string
		id     string
		wantOk bool
	}{
		{"ok", "pk_ok", "psk_ok", "proj_ok", true},
		{"missing key", "", "psk", "proj", false},
		{"missing secret", "pk", "", "proj", false},
		{"missing id", "pk", "psk", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := New(tt.key, tt.secret, tt.id)
			if (err == nil) != tt.wantOk {
				t.Errorf("New() err = %v, wantOk = %v", err, tt.wantOk)
			}
			if tt.wantOk && c == nil {
				t.Error("expected non-nil client")
			}
		})
	}
}

func TestNew_trimSpace(t *testing.T) {
	c, err := New("  pk  ", "  psk  ", "  proj  ")
	if err != nil {
		t.Fatal(err)
	}
	if c.projectKey != "pk" || c.projectSecret != "psk" || c.projectID != "proj" {
		t.Errorf("expected trimmed credentials, got key=%q secret=%q id=%q", c.projectKey, c.projectSecret, c.projectID)
	}
}

func TestSend_injectsProjectID(t *testing.T) {
	var body []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusCreated)
	}))
	defer srv.Close()
	target, _ := url.Parse(srv.URL)
	c, err := New("pk", "psk", "proj_123", WithHTTPClient(&http.Client{
		Transport: redirectTransport{target: target, rt: http.DefaultTransport},
	}))
	if err != nil {
		t.Fatal(err)
	}
	err = c.Send(context.Background(), map[string]any{"message": "hi"})
	if err != nil {
		t.Fatal(err)
	}
	if body == nil {
		t.Fatal("server did not receive body")
	}
	bodyStr := string(body)
	if !strings.Contains(bodyStr, "proj_123") || !strings.Contains(bodyStr, "hi") {
		t.Errorf("body %s missing project_id or message", bodyStr)
	}
}

func TestSend_errorStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"bad key"}`))
	}))
	defer srv.Close()
	target, _ := url.Parse(srv.URL)
	c, _ := New("pk", "psk", "proj", WithHTTPClient(&http.Client{
		Transport: redirectTransport{target: target, rt: http.DefaultTransport},
	}))
	err := c.Send(context.Background(), map[string]any{"message": "x"})
	if err == nil {
		t.Fatal("expected error for 401")
	}
	if !strings.Contains(err.Error(), "401") {
		t.Errorf("error should mention status: %v", err)
	}
}
