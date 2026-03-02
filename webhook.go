package loopengine

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
)

// VerifyWebhook verifies that a webhook payload was signed by LoopEngine.
//
// Pass the raw request body before any JSON unmarshalling — the signature is
// computed over the exact bytes received. Unmarshalling and re-marshalling the
// payload will produce a different byte sequence and break verification.
//
// signatureHeader and timestampHeader are the values of the
// X-LoopEngine-Signature and X-LoopEngine-Timestamp headers respectively.
// If maxAgeSec > 0, the timestamp must be within ±maxAgeSec seconds of now
// (default 300, i.e. 5 minutes). Pass 0 to skip the timestamp check.
//
// Returns true only if the signature matches and the timestamp is within the
// allowed window.
//
// Example (net/http):
//
//	func webhookHandler(w http.ResponseWriter, r *http.Request) {
//	    body, _ := io.ReadAll(r.Body) // read raw bytes BEFORE json.Unmarshal
//	    if !loopengine.VerifyWebhook(
//	        os.Getenv("LOOPENGINE_WEBHOOK_SECRET"),
//	        r.Header.Get("X-LoopEngine-Signature"),
//	        r.Header.Get("X-LoopEngine-Timestamp"),
//	        body,
//	        300,
//	    ) {
//	        http.Error(w, "invalid signature", http.StatusUnauthorized)
//	        return
//	    }
//	    var event map[string]any
//	    json.Unmarshal(body, &event)
//	    // handle event …
//	}
func VerifyWebhook(secret, signatureHeader, timestampHeader string, rawBody []byte, maxAgeSec int) bool {
	if signatureHeader == "" || timestampHeader == "" || len(signatureHeader) < 3 || signatureHeader[:3] != "v1=" {
		return false
	}

	if maxAgeSec > 0 {
		ts, err := strconv.ParseInt(timestampHeader, 10, 64)
		if err != nil {
			return false
		}
		if absInt64(time.Now().Unix()-ts) > int64(maxAgeSec) {
			return false
		}
	}

	// signed content matches server computeSignature: timestamp + "." + rawBody
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(timestampHeader))
	mac.Write([]byte("."))
	mac.Write(rawBody)
	expected := "v1=" + hex.EncodeToString(mac.Sum(nil))

	// hmac.Equal is constant-time: no timing oracle on the comparison.
	return hmac.Equal([]byte(expected), []byte(signatureHeader))
}

func absInt64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
