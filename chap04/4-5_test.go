package chap04

import (
	"log"
	"testing"
)

func Test_checkStatus(t *testing.T) {

	done := make(chan interface{})
	defer close(done)

	urls := []string{"https://www.google.com", "https://badhost"}
	for result := range checkStatus(done, urls...) {
		if result.Error != nil {
			t.Logf("error: %v", result.Error)
			continue
		}
		log.Printf("Response: %v\n", result.Response.Status)
	}
}
