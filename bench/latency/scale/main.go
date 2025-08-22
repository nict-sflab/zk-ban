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

	bb := usefulbench.New(COUNT)
	latency.BenchmarkScalability(bb)

	bb.SaveJson("scale-bench.json")
	fmt.Printf("%v\n", bb.ResultJson())
}
