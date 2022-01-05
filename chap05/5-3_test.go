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


func Test_wrong_heartbeat(t *testing.T) {
	done := make(chan interface{})
	time.AfterFunc(10*time.Second, func() {
		close(done)
	})

	const timeout = 2 * time.Second
	heartbeat, results := doWork(done, timeout/2)
	for {
		select {
		case _, ok := <-heartbeat:
			if ok == false {
				log.Println("heartbeat false")
				return
			}
			log.Println("pulse")
		case r, ok := <-results:
			if ok == false {
				log.Println("results false")
				return
			}
			log.Printf("results %v\n", r.Second())
		case <-time.After(timeout):
			log.Printf("worker goroutine is not healthy!!")
			return
		}
	}

	log.Println("finish")
}
