package chap06

func fib(n int) <-chan int {
	result := make(chan int)
	go func() {
		defer close(result)
		if n <= 2 {
			result <- 1
			return
		}
		result <- <-fib(n-1) + <-fib(n-2)
	}()
	// ここから継続
	// 継続とはプログラム中のある計算処理の途中からその処理を終わらせるまでに行われる処理のまとまりを指す
	return result
}