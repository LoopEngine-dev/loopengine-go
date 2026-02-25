package loopengine_test

import (
	"context"
	"fmt"
	"os"

	"github.com/LoopEngine-dev/loopengine-go"
)

func ExampleClient() {
	client, err := loopengine.New(
		os.Getenv("LOOPENGINE_PROJECT_KEY"),
		os.Getenv("LOOPENGINE_PROJECT_SECRET"),
		os.Getenv("LOOPENGINE_PROJECT_ID"),
	)
	if err != nil {
		fmt.Println("New:", err)
		return
	}
	err = client.Send(context.Background(), map[string]any{"message": "Hello from the SDK"})
	if err != nil {
		fmt.Println("Send:", err)
		return
	}
	fmt.Println("ok")
}
