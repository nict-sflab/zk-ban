package main

import (
	"fmt"

	"github.com/akakou/zk-ban/bench/latency"
	"github.com/akakou/zk-ban/bench/latency/utils/usefulbench"
)

const COUNT = 20

func main() {
	latency.OnlyUniform = false
	bb := usefulbench.New(COUNT)
	latency.BenchmarkScalability(bb)

	bb.SaveJson("scale-bench.json")
	fmt.Printf("%v\n", bb.ResultJson())
}
