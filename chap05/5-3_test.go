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

func Test_heartbeat_per_job(t *testing.T) {
	done := make(chan interface{})
	defer close(done)
	heartbeat, results := doWorkPerJob(done)
	for {
		select {
		case _, ok := <-heartbeat:
			if ok {
				log.Println("pulse")
			} else {
				return
			}
		case r, ok := <-results:
			if ok {
				log.Printf("results: %v\n", r)
			} else {
				return
			}
		}
	}
}

// this test will fail
func Test_doWorkWithIntStream_BadCase(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	intSlice := []int{0, 1, 2, 3, 5}
	_, results := doWorkWithIntStream(done, intSlice...)

	for i, expected := range intSlice {
		select {
		case r := <-results:
			if r != expected {
				t.Errorf(
					"index %v: expected %v, but received %v,",
					i,
					expected,
					r,
				)
			}
		case <-time.After(1 * time.Second):
			//t.Fatalf("test timed out")
			log.Println("test timed out")
			return
		}
	}
}

func Test_doWorkWithIntStream_GoodCase(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	intSlice := []int{0, 1, 2, 3, 5}
	heartbeat, results := doWorkWithIntStream(done, intSlice...)

	<-heartbeat

	for i, expected := range intSlice {
		select {
		case r := <-results:
			if r != expected {
				t.Errorf(
					"index %v: expected %v, but received %v,",
					i,
					expected,
					r,
				)
			}
		}
	}
}

func Test_doWorkMoreSafteyHeartbeat(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	intSlice := []int{0, 1, 2, 3, 5}
	const timeout = 2 * time.Second
	heartbeat, results := doWorkMoreSafeHeartbeat(done, timeout/2, intSlice...)

	<-heartbeat

	i := 0
	for {
		select {
		case r, ok := <-results:
			if ok == false {
				return
			} else if expected := intSlice[i]; r != expected {
				t.Errorf(
					"index %v: expected %v, but received %v,",
					i,
					expected,
					r,
				)
			}
			i++
		case <-heartbeat:
		case <-time.After(timeout):
			t.Fatal("test timed out")
		}
	}
}
