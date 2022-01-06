package addendum

import (
	"context"
	"log"
	"os"
	"runtime/trace"
	"testing"
)

func Test_trace(t *testing.T) {
	f, err := os.Create("trace.out")
	if err != nil {
		log.Fatalf("failed to create trace: %v", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatalf("failed to close trace: %v", err)
		}
	}()
	if err := trace.Start(f); err != nil {
		panic(err)
	}
	defer trace.Stop()

	ctx := context.Background()
	ctx, task := trace.NewTask(ctx, "makeCoffee")
	defer task.End()
	trace.Log(ctx, "orderID", "1")

	coffee := make(chan bool)

	go func() {
		trace.WithRegion(ctx, "extractCoffee", func() {

		})
		coffee <-true
	}()
	<-coffee
}
