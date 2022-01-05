package chap05

import (
	"go-concurrency/chap04"
	"log"
	"time"
)

type startGoroutineFn func(
	done <-chan interface{},
	pulseInterval time.Duration,
) (heartbeat <-chan interface{})

// chap4から引用
func or(channels ...<-chan interface{}) <-chan interface{} {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	orDone := make(chan interface{})
	go func() {
		defer close(orDone)

		switch len(channels) {
		case 2:
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default:
			select {
			case <-channels[0]:
			case <-channels[1]:
			case <-channels[2]:
			case <-or(append(channels[3:], orDone)...):
			}
		}
	}()
	return orDone
}

func newSteward(
	timeout time.Duration,
	startGoroutine startGoroutineFn,
) startGoroutineFn {
	return func(
		done <-chan interface{},
		pulseInterval time.Duration,
	) <-chan interface{} {
		heartbeat := make(chan interface{})
		go func() {
			defer close(heartbeat)

			var wardDone chan interface{}
			var wardHeartbeat <-chan interface{}
			startWard := func() {
				wardDone = make(chan interface{})
				wardHeartbeat = startGoroutine(or(wardDone, done), timeout/2)
			}
			startWard()
			pulse := time.Tick(pulseInterval)

		monitorLoop:
			for {
				timeoutSignal := time.After(timeout)

				for {
					select {
					case <-pulse:
						select {
						case heartbeat <- struct{}{}:
						default:
						}
					case <-wardHeartbeat:
						continue monitorLoop
					case <-timeoutSignal:
						log.Println("steward: ward unhealthy; restarting")
						close(wardDone)
						startWard()
						continue monitorLoop
					case <-done:
						return
					}
				}
			}
		}()

		return heartbeat
	}
}

func doWorkFn(
	done <-chan interface{},
	intList ...int,
) (startGoroutineFn, <-chan interface{}) {
	intChanStream := make(chan (<-chan interface{}))
	intStream := chap04.Bridge(done, intChanStream)
	doWork := func(done <-chan interface{}, pulseInterval time.Duration) <-chan interface{} {
		intStream := make(chan interface{})
		heartbeat := make(chan interface{})

		go func() {
			defer close(intStream)
			select {
			case intChanStream <- intStream:
			case <-done:
				return
			}
			pulse := time.Tick(pulseInterval)

			for {
			valueLoop:
				for _, intVal := range intList {
					if intVal < 0 {
						log.Printf("negative value: %v\n", intVal)
						return
					}

					for {
						select {
						case <-pulse:
							select {
							case heartbeat <- struct{}{}:
							default:
							}
						case intStream <- intVal:
							continue valueLoop
						case <-done:
							return
						}
					}
				}
			}
		}()
		return heartbeat
	}

	return doWork, intStream
}

func take(done <-chan interface{}, valueStream <-chan interface{}, num int) <-chan interface{} {
	takeStream := make(chan interface{})
	go func() {
		defer close(takeStream)
		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			case takeStream <- <-valueStream:
			}
		}
	}()
	return takeStream
}