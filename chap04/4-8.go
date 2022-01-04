package chap04

func orDone(done, c <-chan interface{}) <-chan interface{} {
	valStream := make(chan interface{})
	go func() {
		defer close(valStream)
		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				if ok == false {
					return
				}
				select {
				case valStream<-v:
				case <-done:
				}
			}
		}
	}()
	return valStream
}

// Usecase
// for val := range orDone(done, myChan) { // val に対して何かする　}
// myChanが外部から取得する場合に、複雑なエラー制御をorDone内に閉じ込めてしまう。