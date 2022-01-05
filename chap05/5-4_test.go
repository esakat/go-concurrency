package chap05

import (
	"log"
	"sync"
	"testing"
)

func Test_doWorkCopy(t *testing.T) {
	done := make(chan interface{})
	result := make(chan int)

	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go doWorkCopy(done, i, &wg, result)
	}

	firstReturned := <-result
	close(done)
	wg.Wait()

	log.Printf("Received an answer from %v\n", firstReturned)
}
