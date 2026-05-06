package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/akakou/zk-ban/bench/storage"
	"github.com/akakou/zk-ban/test"
	"github.com/akakou/zk-ban/witness"
)

const KAPPA = 30
const LAMBDA = 32768

func main() {
	result := make(storage.Result, 0)

	rlSize := witness.MakeUniformRLSizeFromTotal(KAPPA, LAMBDA)
	storage.BenchUpdate(rlSize, "compare-circuit", "uniform", "compare-storage", result)

	j, err := json.Marshal(result)

	test.PanicIfErr(err)
	fmt.Printf("%s", j)

	os.WriteFile("compare-storage.json", j, 0o644)
}
