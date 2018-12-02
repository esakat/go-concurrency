package main

import (
	"fmt"
)

func main() {
	intStream := make(chan interface{})
	go func() {
		defer close(intStream)
		for i := 1; i <= 5; i++ {
			intStream <- i
			intStream <- "a"
		}
	}()

	for integer := range intStream {
		fmt.Printf("%v ", integer)
	}

	ch := make(chan int, 4)
	ch <- 1
	ch <- 1
	ch <- 1
	ch <- 1
	ch <- 1
	fmt.Printf("%v ", <-ch)
}
