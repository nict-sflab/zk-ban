package main

import (
	"fmt"

	gnarkserializable "github.com/akakou/gnark-serializable"
	"github.com/akakou/zk-ban/bench/latency"
	"github.com/akakou/zk-ban/bench/latency/utils/usefulbench"
)

const COUNT = 20

func main() {
	gnarkserializable.Unsafe = true
	latency.NoParallel = true

	bb := usefulbench.New(COUNT)
	latency.BenchmarkBaseline(bb)

	bb.SaveJson("compare-bench.json")
	fmt.Printf("%v\n", bb.ResultJson())
}
