package main

import (
	"fmt"

	"github.com/akakou/zk-ban/bench/latency"
	"github.com/akakou/zk-ban/bench/latency/utils/usefulbench"
)

const COUNT = 20

func main() {

	bb := usefulbench.New(COUNT)
	latency.BenchmarkBaseline(bb)

	bb.SaveJson("baseline-bench.json")
	fmt.Printf("%v\n", bb.ResultJson())
}
