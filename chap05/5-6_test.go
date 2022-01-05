package chap05

import (
	"log"
	"os"
	"testing"
	"time"
)

func Test_Monitor(t *testing.T) {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	doWork := func(done <-chan interface{}, _ time.Duration) <-chan interface{} {
		log.Println("ward: Hello, Im irresponsible")
		go func() {
			<-done
			log.Println("ward: I am halting.")
		}()
		return nil
	}

	doWorkWitSteward := newSteward(4*time.Second, doWork)

	done := make(chan interface{})
	time.AfterFunc(9*time.Second, func() {
		log.Println("main: halting steward and ward.")
		close(done)
	})

	for range doWorkWitSteward(done, 4*time.Second) {
	}
	log.Println("Done!")
}

func Test_Monitor2(t *testing.T) {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	done := make(chan interface{})
	defer close(done)

	doWork, intStream := doWorkFn(done, 1, 2, -1, 3, 4, 5)
	doWorkWithSteward := newSteward(1*time.Millisecond, doWork)
	doWorkWithSteward(done, 1*time.Hour)


	for intVal := range take(done, intStream, 6) {
		log.Printf("Received: %v\n", intVal)
	}
	log.Println("Done!")
}
