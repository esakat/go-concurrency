package main

import (
	"fmt"
	"sync"
)

func main() {
	type Button struct {
		Clicked *sync.Cond
	}

	button := Button{Clicked: sync.NewCond(&sync.Mutex{})}

	subscribe := func(c *sync.Cond, fn func()) {
		var goroutineRunning sync.WaitGroup
		goroutineRunning.Add(1)
		go func() {
			goroutineRunning.Done()
			c.L.Lock()
			defer c.L.Unlock()
			c.Wait()
			fn()
		}()
		goroutineRunning.Wait()
	}

	var clickRegisterd sync.WaitGroup
	clickRegisterd.Add(3)
	subscribe(button.Clicked, func() {
		fmt.Println("Maximizing Window.")
		clickRegisterd.Done()
	})
	subscribe(button.Clicked, func() {
		fmt.Println("Open Box!!")
		clickRegisterd.Done()
	})
	subscribe(button.Clicked, func() {
		fmt.Println("Mouse Clicked!")
		clickRegisterd.Done()
	})

	button.Clicked.Broadcast()

	clickRegisterd.Wait()
}
