package chap04

import (
	"log"
	"testing"
)

func Test_Tee(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	out1, out2 := tee(done, take(done, repeat(done, 1, 2), 4))

	for val1 := range out1 {
		log.Printf("out1: %v, out2: %v\n", val1, <-out2)
	}
}
