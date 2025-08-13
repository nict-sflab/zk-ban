package main

import (
	"fmt"

	"github.com/akakou/zk-ban/test"
	"github.com/akakou/zk-ban/test/utils/usefulbench"
)

const COUNT = 20

func main() {
	bb := usefulbench.New(COUNT)
	test.BenchmarkUpdate(bb)

	bb.SaveJson("update-bench.json")
	fmt.Printf("%v\n", bb.ResultJson())
}
