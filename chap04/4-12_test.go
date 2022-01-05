package chap04

import (
	"context"
	"log"
	"sync"
	"testing"
)

func Test_done_style(t *testing.T) {
	var wg sync.WaitGroup
	done := make(chan interface{})
	defer close(done)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printGreeting(done); err != nil {
			log.Printf("%v", err)
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printFarewell(done); err != nil {
			log.Printf("%v", err)
			return
		}
	}()

	wg.Wait()
}

func Test_ctx_style(t *testing.T) {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printGreetingWithCtx(ctx); err != nil {
			log.Printf("cannot print greeting: %v\n", err)
			cancel()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printFarewellWithCtx(ctx); err != nil {
			log.Printf("cannot print farewell: %v", err)
		}
	}()

	wg.Wait()
}

func Test_ctx_with_val(t *testing.T) {
	ProcessRequest("jane", "abc123")
}