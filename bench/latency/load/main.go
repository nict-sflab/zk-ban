package main

import (
	"fmt"

	"github.com/akakou/zk-ban/bench/latency"
	"github.com/akakou/zk-ban/bench/latency/utils/usefulbench"
	"github.com/akakou/zk-ban/load"
)

const COUNT = 20

func main() {
	load.KeyChecked = true

	bb := usefulbench.New(COUNT)
	latency.BenchmarkLoad(bb)

	bb.SaveJson("load-bench.json")
	fmt.Printf("%v\n", bb.ResultJson())
}
