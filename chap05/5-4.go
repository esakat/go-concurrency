package chap05

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

func doWorkCopy(
	done <-chan interface{},
	id int,
	wg *sync.WaitGroup,
	result chan<- int,
	) {
	started := time.Now()
	defer wg.Done()

	simulatedLoadTime := time.Duration(1+rand.Intn(5)) * time.Second
	select {
	case <-done:
	case <-time.After(simulatedLoadTime):
	}

	select {
	case <-done:
	case result <- id:
	}

	took := time.Since(started)
	if took < simulatedLoadTime {
		took = simulatedLoadTime
	}
	log.Printf("%v took %v\n", id, took)
}