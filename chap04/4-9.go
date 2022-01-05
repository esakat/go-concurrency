package chap04

func tee(done <-chan interface{}, in <-chan interface{}) (_, _ <-chan interface{}) {
	out1 := make(chan interface{})
	out2 := make(chan interface{})
	go func() {
		defer close(out1)
		defer close(out2)
		for val := range orDone(done, in) {
			// ローカル変数定義
			var out1, out2 = out1, out2
			// 確実に両方に書き込めるように2回ループさせる
			for i := 0; i < 2; i++ {
				select {
				// 片方に書き込みが成功したら
				case out1 <- val:
					// out1をnilで上書きして、もう使えなくする(次のループでは1回目で使わなかった方に確実に書き込まれる
					out1 = nil
				case out2 <- val:
					out2 = nil
				}
			}
		}
	}()
	return out1, out2
}
