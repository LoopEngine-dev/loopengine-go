package loopengine

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"
	"time"
)

const testSecret = "whsec_live_test_secret"

var testBody = []byte(`{"event":"feedback.created","id":"evt_123"}`)

// makeSignature generates a valid (timestamp, "v1=<hex>") pair for the given inputs.
func makeSignature(secret string, body []byte, ts string) (string, string) {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(ts))
	mac.Write([]byte("."))
	mac.Write(body)
	return ts, "v1=" + hex.EncodeToString(mac.Sum(nil))
}

func nowTS() string {
	return fmt.Sprintf("%d", time.Now().Unix())
}

func TestVerifyWebhook_valid(t *testing.T) {
	ts, sig := makeSignature(testSecret, testBody, nowTS())
	if !VerifyWebhook(testSecret, sig, ts, testBody, 300) {
		t.Error("expected true for valid signature")
	}
}

func TestVerifyWebhook_table(t *testing.T) {
	ts, validSig := makeSignature(testSecret, testBody, nowTS())
	oldTS := fmt.Sprintf("%d", time.Now().Unix()-600)
	_, oldSig := makeSignature(testSecret, testBody, oldTS)
	altered := []byte(`{"event":"feedback.created","id":"evt_456"}`)
	_, wrongSig := makeSignature("wrong_secret", testBody, ts)
	tampered := validSig[:len(validSig)-4] + "aaaa"

	tests := []struct {
		name      string
		secret    string
		sig       string
		tsHeader  string
		body      []byte
		maxAge    int
		wantValid bool
	}{
		{"valid", testSecret, validSig, ts, testBody, 300, true},
		{"tampered signature", testSecret, tampered, ts, testBody, 300, false},
		{"wrong secret", testSecret, wrongSig, ts, testBody, 300, false},
		{"altered body", testSecret, validSig, ts, altered, 300, false},
		{"replay old ts", testSecret, oldSig, oldTS, testBody, 300, false},
		{"replay disabled (maxAge=0)", testSecret, oldSig, oldTS, testBody, 0, true},
		{"missing v1= prefix", testSecret, validSig[3:], ts, testBody, 300, false},
		{"empty signature header", testSecret, "", ts, testBody, 300, false},
		{"empty timestamp header", testSecret, validSig, "", testBody, 300, false},
		{"non-numeric timestamp", testSecret, validSig, "not-a-number", testBody, 300, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := VerifyWebhook(tc.secret, tc.sig, tc.tsHeader, tc.body, tc.maxAge)
			if got != tc.wantValid {
				t.Errorf("VerifyWebhook() = %v, want %v", got, tc.wantValid)
			}
		})
	}
}
