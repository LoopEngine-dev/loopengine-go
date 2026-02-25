package loopengine

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"strconv"
	"strings"
	"time"
)

// signRequest returns timestamp and signature for the given method, path, and body.
// Signature is HMAC-SHA256 over "METHOD\nPATH\nTIMESTAMP\nSHA256(body)", base64url.
func signRequest(secret, method, path string, body []byte) (timestamp, signature string) {
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	sum := sha256.Sum256(body)
	bodyHash := hex.EncodeToString(sum[:])
	canonical := strings.Join([]string{method, path, ts, bodyHash}, "\n")
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(canonical))
	sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return ts, "v1=" + sig
}
