package chap04

import (
	"log"
	"testing"
)

func Test_Bridge(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	for v := range bridge(nil, genVals()) {
		log.Printf("%v ", v)
	}
}
