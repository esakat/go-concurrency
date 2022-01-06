package chap06

import (
	"log"
	"testing"
)

func Test_fib(t *testing.T) {
	log.Printf("fib(4) = %d", <-fib(4))
}
