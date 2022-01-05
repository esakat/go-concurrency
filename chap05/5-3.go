package chap05

import (
	"math/rand"
	"time"
)

// heartbeat

func doWork(
	done <-chan interface{},
	pulseInterval time.Duration,
) (<-chan interface{}, <-chan time.Time) {
	heartbeat := make(chan interface{})
	results := make(chan time.Time)
	go func() {
		defer close(heartbeat)
		defer close(results)

		pulse := time.Tick(pulseInterval)
		workGen := time.Tick(2 * pulseInterval)

		sendPulse := func() {
			select {
			case heartbeat <- struct{}{}:
			default:
			}
		}

		sendResult := func(r time.Time) {
			for {
				select {
				case <-done:
					return
				case <-pulse:
					sendPulse()
				case results <- r:
					return
				}
			}
		}

		for {
			select {
			case <-done:
				return
			case <-pulse:
				sendPulse()
			case r := <-workGen:
				sendResult(r)
			}
		}

	}()
	return heartbeat, results
}

func doWorkPerJob(done <-chan interface{}) (<-chan interface{}, <-chan int) {
	heartbeatStream := make(chan interface{}, 1)
	workStream := make(chan int)
	go func() {
		defer close(heartbeatStream)
		defer close(workStream)

		for i := 0; i < 10; i++ {
			select {
			case heartbeatStream <- struct{}{}:
			default:
			}

			select {
			case <-done:
				return
			case workStream <- rand.Intn(10):
			}
		}
	}()

	return heartbeatStream, workStream
}


func doWorkWithIntStream(
	done <-chan interface{},
	nums ...int,
	) (<-chan interface{}, <-chan int) {
	heartbeatStream := make(chan interface{}, 1)
	intStream := make(chan int)
	go func() {
		defer close(heartbeatStream)
		defer close(intStream)

		time.Sleep(2*time.Second)

		for _, n := range nums {
			select {
			case heartbeatStream <- struct{}{}:
			default:
			}

			select {
			case <-done:
				return
			case intStream <- n:
			}
		}
	}()

	return heartbeatStream, intStream
}



func doWorkMoreSafeHeartbeat(
	done <-chan interface{},
	pulseInterval time.Duration,
	nums ...int,
) (<-chan interface{}, <-chan int) {
	heartbeat := make(chan interface{}, 1)
	intStream := make(chan int)
	go func() {
		defer close(heartbeat)
		defer close(intStream)

		time.Sleep(2*time.Second)

		pulse := time.Tick(pulseInterval)

		numLoop:
		for _, n := range nums {
			for {
				select {
				case <-done:
					return
				case <-pulse:
					select {
					case heartbeat <- struct{}{}:
					default:
					}
				case intStream <- n:
					continue numLoop
				}
			}
		}
	}()

	return heartbeat, intStream
}
