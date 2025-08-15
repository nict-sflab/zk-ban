package usefulbench_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/akakou/zk-ban/bench/latency/utils/usefulbench"
)

func TestBench(t *testing.B) {
	b := usefulbench.New(3)

	oneSecond := time.Second * 1
	twoSecond := time.Second * 2
	threeSecond := time.Second * 3
	fourSecond := time.Second * 4

	b.Run("test:1", func(b usefulbench.Benchmarker) {
		for b.Loop() {
			time.Sleep(twoSecond)
		}
	})

	b.Run("test:2", func(b usefulbench.Benchmarker) {
		for b.Loop() {
			time.Sleep(threeSecond)
		}
	})

	firstResult := b.Result["test"]["1"]
	secondResult := b.Result["test"]["2"]

	if firstResult < int64(oneSecond) || firstResult > int64(threeSecond) {
		t.Fatalf("First Result is wrong: not match to %d(1s) < %d(firstResult) < %d(3s)", oneSecond, firstResult, threeSecond)
	}

	if secondResult < int64(twoSecond) || secondResult > int64(fourSecond) {
		t.Fatalf("Second Result is wrong: not match to %d(2s) < %d(secondResult) < %d(4s)", twoSecond, firstResult, fourSecond)
	}

	fmt.Printf("%s\n", b.ResultJson())
	fmt.Printf("first: %v, second : %v\n", firstResult, secondResult)

	panic("")
}
