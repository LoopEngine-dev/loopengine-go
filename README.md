# LoopEngine Go SDK

Go client for the [LoopEngine](https://loopengine.dev) Ingest API. Create a client with your credentials, then call `Send` with your payload.

**Requirements:** Go 1.21+

## Installation

```bash
go get github.com/LoopEngine-dev/loopengine-go
```

## Usage

```go
package main

import (
    "context"
    "log"

    "github.com/LoopEngine-dev/loopengine-go"
)

func main() {
    client, err := loopengine.New(projectKey, projectSecret, projectID)
    if err != nil {
        log.Fatal(err)
    }
    err = client.Send(context.Background(), map[string]any{"message": "User feedback here"})
    if err != nil {
        log.Fatal(err)
    }
}
```

- **New(key, secret, projectID)** — Builds a client. Use your project key, secret, and project ID from the LoopEngine dashboard.
- **Send(ctx, payload)** — Sends the payload to the Ingest API at `api.loopengine.dev`. `project_id` is added automatically. The payload must match the **fields and constraints** configured for your project in the LoopEngine dashboard (e.g. required fields, allowed keys, value types). You can pass a `map[string]any`, a struct, or any JSON-serializable value that conforms to your project’s schema.

The client is safe for concurrent use and has no dependencies beyond the standard library. You can pass **WithHTTPClient** to set timeouts or a custom transport:

```go
client, err := loopengine.New(key, secret, projectID,
    loopengine.WithHTTPClient(&http.Client{Timeout: 10 * time.Second}),
)
```

## Geolocation

You can send device location so feedback is associated with coordinates instead of IP-based geo. Pass an optional third argument to `Send` with `SendOptions` containing both `GeoLat` and `GeoLon`. When both are non-zero, the SDK adds `geo_lat` and `geo_lon` to the request body; they are included in the HMAC signature. Omit the third argument (or pass `nil`) to use IP-based geolocation. Valid ranges: latitude -90 to 90, longitude -180 to 180.

```go
// Without geo (IP-based location is used)
err = client.Send(ctx, map[string]any{"message": "Feedback"})

// With device coordinates
err = client.Send(ctx, map[string]any{"message": "Bug at my location"},
    &loopengine.SendOptions{GeoLat: 34.05, GeoLon: -118.25})
```

## Verifying webhook payloads

When LoopEngine delivers a webhook to your endpoint, it signs the request with HMAC-SHA256 using a **signing secret** that only you and LoopEngine know. Verifying that signature before processing the event confirms the request came from LoopEngine, was not tampered with, and — by also checking the timestamp — limits replay attacks.

**Get the secret:** In your dashboard, open your project → Webhooks. The signing secret (`whsec_live_...`) is shown when you create or rotate the webhook. Store it as an environment variable and never commit it.

**Critical:** call `VerifyWebhook` with the raw request body bytes **before** any JSON unmarshalling. The signature is computed over the exact bytes received; unmarshalling and re-marshalling the payload will produce a different byte sequence and fail verification.

### Signature

```go
func loopengine.VerifyWebhook(
    secret           string,
    signatureHeader  string,
    timestampHeader  string,
    rawBody          []byte,
    maxAgeSec        int,
) bool
```

### Parameters

| Parameter | Type | Description |
|---|---|---|
| `secret` | `string` | Signing secret from the dashboard (`whsec_live_...`) |
| `signatureHeader` | `string` | Full value of the `X-LoopEngine-Signature` header |
| `timestampHeader` | `string` | Value of the `X-LoopEngine-Timestamp` header (Unix seconds) |
| `rawBody` | `[]byte` | Raw HTTP body as received, before any JSON unmarshalling |
| `maxAgeSec` | `int` | Max timestamp age in seconds; use `300` for 5 min. Pass `0` to skip. |

**Returns:** `bool` — `true` if valid, `false` if the signature does not match or the timestamp is outside the allowed window.

### net/http example

```go
package main

import (
    "encoding/json"
    "io"
    "net/http"
    "os"

    loopengine "github.com/LoopEngine-dev/loopengine-go"
)

func webhookHandler(w http.ResponseWriter, r *http.Request) {
    body, err := io.ReadAll(r.Body) // read raw bytes BEFORE json.Unmarshal
    if err != nil {
        http.Error(w, "cannot read body", http.StatusBadRequest)
        return
    }

    if !loopengine.VerifyWebhook(
        os.Getenv("LOOPENGINE_WEBHOOK_SECRET"),
        r.Header.Get("X-LoopEngine-Signature"),
        r.Header.Get("X-LoopEngine-Timestamp"),
        body,
        300,
    ) {
        http.Error(w, "invalid signature", http.StatusUnauthorized)
        return
    }

    var event map[string]any
    if err := json.Unmarshal(body, &event); err != nil { // parse AFTER verifying
        http.Error(w, "invalid JSON", http.StatusBadRequest)
        return
    }

    // handle event …
    w.WriteHeader(http.StatusOK)
}
```
