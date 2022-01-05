package chap05

import (
	"log"
	"testing"
	"time"
)

func Test_heartbeat(t *testing.T) {
	done := make(chan interface{})
	time.AfterFunc(10*time.Second, func() {
		close(done)
	})

	const timeout = 2 * time.Second
	heartbeat, results := doWork(done, timeout/2)
	for {
		select {
		case _, ok := <-heartbeat:
			if !ok {
				return
			}
			log.Println("pulse")
		case r, ok := <-results:
			if ok == false {
				return
			}
			log.Printf("results %v\n", r.Second())
		case <-time.After(timeout):
			return
		}
	}
}
