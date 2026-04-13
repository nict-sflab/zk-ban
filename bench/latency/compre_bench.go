package latency

import (
	"github.com/akakou/zk-ban/bench/latency/utils/usefulbench"
	"github.com/akakou/zk-ban/test"
)

func BenchmarkCompare(b usefulbench.Benchmarker) {
	params := test.PrepareParams()

	baseNym := 1024*32 + 16*14
	basePeriod := 30

	ub, ok := b.(*usefulbench.UsefulBenchmaker)

	if ok {
		ub.Result["env"]["basePeriodNum"] = int64(basePeriod)
		ub.Result["env"]["baseNymNum"] = int64(baseNym)
	}

	benchmarkBaseline(&params, b)
	rlU := test.EmptyUniformRevocationList(basePeriod, baseNym)
	benchmarkBasicUpdate(rlU, "compare", &params, b)
}
