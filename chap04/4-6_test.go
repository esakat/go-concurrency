package chap04

import (
	"log"
	"testing"
)

// パイプラインに組み込める条件
// 1. 関数は引数と返り値が同じ型であること
// パイプラインはバッチ処理に近い(一気に全ての項目を次の関数に渡してる)
// -> ストリーム処理の場合、1つづつ渡されるらしい
// パイプラインは一時的に引数の項目数の倍のメモリふっとプリントが取られちゃうのが欠点
func Test_pipeline(t *testing.T) {
	ints := []int{1, 2, 3, 4}
	for _, v := range _add(_multiply(ints, 2), 1) {
		log.Println(v)
	}
}

func Test_generator(t *testing.T) {
	done := make(chan interface{})
	defer close(done)
	intStream := generator(done, 1, 2, 3, 4)
	pipeline := multiply(done, add(done, multiply(done, intStream, 2), 1), 2)
	for v := range pipeline {
		log.Println(v)
	}
}
