package loopengine

import (
	"testing"
)

func TestSignRequest_deterministic(t *testing.T) {
	secret := "psk_test"
	method := "POST"
	path := "/feedback"
	body := []byte(`{"project_id":"proj","message":"hi"}`)
	ts1, sig1 := signRequest(secret, method, path, body)
	ts2, sig2 := signRequest(secret, method, path, body)
	// Timestamps may differ by 1 second if we cross a second boundary
	if sig1 != sig2 && ts1 == ts2 {
		t.Error("same inputs should produce same signature when timestamp equal")
	}
	if len(ts1) == 0 || len(sig1) == 0 {
		t.Error("timestamp and signature must be non-empty")
	}
	if sig1[:3] != "v1=" {
		t.Errorf("signature should have v1= prefix, got %q", sig1)
	}
}
