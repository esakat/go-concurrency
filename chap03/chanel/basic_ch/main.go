package main

import (
	"fmt"
)

func main() {
	stringStream := make(chan string)
	go func() {
		stringStream <- "Hello Stream!!"
		defer close(stringStream)
	}()
	text, ok := <-stringStream
	fmt.Printf("(%v): %v", ok, text)
}
