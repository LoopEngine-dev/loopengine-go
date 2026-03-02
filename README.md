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
