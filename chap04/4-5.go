package chap04

import "net/http"

// Errorと結果を対にすることが大切
// もしゴルーチンがエラーを発生させる可能性があるのであれば、正常系と強く結びつけて第一級市民として扱う必要がある
type Result struct {
	Error    error
	Response *http.Response
}

func checkStatus(done <-chan interface{}, urls ...string) <-chan Result {
	results := make(chan Result)
	go func() {
		defer close(results)
		for _, url := range urls {
			var result Result
			resp, err := http.Get(url)
			result = Result{Error: err, Response: resp}
			select {
			case <-done:
				return
			case results <- result:
			}
		}
	}()
	return results
}
