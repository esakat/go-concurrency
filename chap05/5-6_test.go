package chap05

import (
	"log"
	"os"
	"testing"
	"time"
)

func Test_Monitor(t *testing.T) {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime|log.LUTC)

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

	for range doWorkWitSteward(done, 4*time.Second) {}
	log.Println("Done!")
}
