// Package loopengine provides a minimal, high-performance client for the
// LoopEngine Ingest API. Create a client with your project credentials, then
// call Send with your payload.
//
//	client, err := loopengine.New(projectKey, projectSecret, projectID)
//	if err != nil { ... }
//	err = client.Send(ctx, map[string]any{"message": "user feedback"})
//
// All signing and HTTP details are handled inside the package. The client is
// safe for concurrent use and uses minimal allocations.
package loopengine
