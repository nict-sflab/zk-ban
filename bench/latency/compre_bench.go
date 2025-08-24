package latency

import (
	"github.com/akakou/zk-ban/bench/latency/utils/usefulbench"
	"github.com/akakou/zk-ban/test"
)

func BenchmarkCompare(b usefulbench.Benchmarker) {
	params := test.PrepareParams()

	baseNym := 114688
	baseSess := 60

	ub, ok := b.(*usefulbench.UsefulBenchmaker)

	if ok {
		ub.Result["env"]["baseSessionNum"] = int64(baseSess)
		ub.Result["env"]["baseNymNum"] = int64(baseNym)
	}

	benchmarkBaseline(&params, b)
	rlU := test.EmptyUniformRevocationList(baseSess, baseNym)

	NoParallel = true
	benchmarkBasicUpdate(rlU, "compare", &params, b)
}
