// Package loopengine provides a minimal, high-performance client for the
// LoopEngine Ingest API. Create a client with your project credentials, then
// call Send with your payload. Optionally pass SendOptions with GeoLat/GeoLon
// to send device coordinates (included in the signed body).
//
//	client, err := loopengine.New(projectKey, projectSecret, projectID)
//	if err != nil { ... }
//	err = client.Send(ctx, map[string]any{"message": "user feedback"})
//	// With geo: err = client.Send(ctx, payload, &loopengine.SendOptions{GeoLat: 34, GeoLon: -118})
//
// All signing and HTTP details are handled inside the package. The client is
// safe for concurrent use and uses minimal allocations.
package loopengine
