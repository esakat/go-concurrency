package main

import (
	"sync"
	"testing"
)

func BenchmarkContextSwitch(b *testing.B) {
	var wg sync.WaitGroup
	begin := make(chan struct{})
	c := make(chan struct{})

	var token struct{}
	sender := func() {
		defer wg.Done()
		<-begin // ❶
		for i := 0; i < b.N; i++ {
			c <- token // ❷
		}
	}
	receiver := func() {
		defer wg.Done()
		<-begin // ❶
		for i := 0; i < b.N; i++ {
			<-c // ❸
		}
	}

	wg.Add(2)
	go sender()
	go receiver()
	b.StartTimer() // ❹
	close(begin)   // ❺
	wg.Wait()
}
