package chap04

// 4.6.1
func _multiply(values []int, multiplier int) []int {
	multipliedValues := make([]int, len(values))
	for i, v := range values {
		multipliedValues[i] = v * multiplier
	}
	return multipliedValues
}

// 4.6.1
func _add(values []int, additive int) []int {
	addedValues := make([]int, len(values))
	for i, v := range values {
		addedValues[i] = v + additive
	}
	return addedValues
}

// 4.6.2
// 個別の数値の塊をチャネル上を流れるストリームに変換する関数
func generator(done <-chan interface{}, integers ...int) <-chan int {
	intStream := make(chan int, len(integers))
	go func() {
		defer close(intStream)
		for _, i := range integers {
			select {
			case <-done:
				return
			// 作成したチャネルに引数で受け取った数値を1つづつ渡していってる
			case intStream <- i:
			}
		}
	}()
	return intStream
}

func multiply(
	done <-chan interface{},
	intStream <-chan int,
	multiplier int,
) <-chan int {
	multipliedStream := make(chan int)
	go func() {
		defer close(multipliedStream)
		for i := range intStream {
			select {
			case <-done:
				return
			case multipliedStream <- i * multiplier:
			}
		}
	}()
	return multipliedStream
}

func add(
	done <-chan interface{},
	intStream <-chan int,
	additive int,
) <-chan int {
	addedStream := make(chan int)
	go func() {
		defer close(addedStream)
		for i := range intStream {
			select {
			case <-done:
				return
			case addedStream <- i + additive:
			}
		}
	}()
	return addedStream
}
