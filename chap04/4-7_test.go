package chap04

import (
	"log"
	"math/rand"
	"runtime"
	"testing"
	"time"
)

// 実行遅い素数を探す例
func Test_SlowSearchPrimeNumber(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	rand := func() interface{} { return rand.Intn(500000000) }

	start := time.Now()

	randIntStream := toInt(done, repeatFn(done, rand))
	log.Println("Primes:")

	for prime := range take(done, primeFinder(done, randIntStream), 10) {
		log.Printf("\t%d\n", prime)
	}
	log.Printf("Search took: %v", time.Since(start))
}

func Test_FanIn(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	start := time.Now()

	rand := func() interface{} { return rand.Intn(500000000) }
	randIntStream := toInt(done, repeatFn(done, rand))

	numFinders := runtime.NumCPU()
	log.Printf("Spinning up %d prime finders.\n", numFinders)
	finders := make([]<-chan interface{}, numFinders)
	log.Println("Primes:")

	for i := 0; i < numFinders; i++ {
		finders[i] = primeFinder(done, randIntStream)
	}
	for prime := range take(done, fanIn(done, finders...), 10) {
		log.Printf("\t%d\n", prime)
	}
	log.Printf("Search took: %v", time.Since(start))
}
