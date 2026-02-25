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
