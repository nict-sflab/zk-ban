package latency

import (
	"fmt"
	"runtime"

	"github.com/akakou/zk-ban/bench"
	"github.com/akakou/zk-ban/bench/latency/utils/usefulbench"
	"github.com/akakou/zk-ban/test"
	zkbanw "github.com/akakou/zk-ban/witness"
	// 	groth16_bls12381 "github.com/consensys/gnark/backend/groth16/bls12-381"
	// fr_bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
)

func BenchmarkLoad(b usefulbench.Benchmarker) {
	BenchmarkNymLoad(b)
	BenchSessLoad(b)
}

func BenchmarkNymLoad(b usefulbench.Benchmarker) {
	logIfUsefulBenchInScaleBench(b)

	// increase nym
	for i := 1; i <= max; i++ {
		runtime.GC()

		nymNum := baseNymNum * i * alpha
		rl := test.EmptyUniformRevocationList(baseSessionNum, nymNum)
		benchmarkLoad(nymNum, rl, bench.NYM_INCREASE_UNIFORM, b)
	}

	// for i := 1; i <= max; i++ {
	// 	runtime.GC()

	// 	nymNum := baseNymNum * i * alpha
	// 	rl := test.EmptyProportionalRevocationList(baseSessionNum, nymNum)
	// 	benchmarkLoad(nymNum, rl, bench.NYM_INCREASE_PROPORTIONAL, b)
	// }

	// for i := 1; i <= max; i++ {
	// 	runtime.GC()

	// 	nymNum := baseNymNum * i * alpha
	// 	rl := test.EmptyGaussianRevocationList(baseSessionNum, nymNum)
	// 	benchmarkLoad(nymNum, rl, bench.NYM_INCREASE_GAUSSIAN, b)
	// }
}

func BenchSessLoad(b usefulbench.Benchmarker) {
	// increase sessionNumber
	BenchSessUniformLoad(b)
	// BenchSessProportionalLoad(b)
	// BenchSessGaussLoad(b)

}

func BenchSessUniformLoad(b usefulbench.Benchmarker) {
	logIfUsefulBenchInScaleBench(b)

	for i := 1; i <= max; i++ {
		runtime.GC()

		sessionNum := baseSessionNum * i * beta
		rl := test.EmptyUniformRevocationList(sessionNum, baseNymNum)
		benchmarkLoad(sessionNum, rl, bench.SESS_INCREASE_UNIFORM, b)
	}

}

func BenchSessProportionalLoad(b usefulbench.Benchmarker) {
	logIfUsefulBenchInScaleBench(b)

	for i := 1; i <= max; i++ {
		runtime.GC()

		sessionNum := baseSessionNum * i * beta
		rl := test.EmptyProportionalRevocationList(sessionNum, baseNymNum)
		benchmarkLoad(sessionNum, rl, bench.SESS_INCREASE_PROPORTIONAL, b)
	}
}

func BenchSessGaussLoad(b usefulbench.Benchmarker) {
	logIfUsefulBenchInScaleBench(b)

	for i := 1; i <= max; i++ {
		runtime.GC()

		sessionNum := baseSessionNum * i * beta
		rl := test.EmptyGaussianRevocationList(sessionNum, baseNymNum)
		benchmarkLoad(sessionNum, rl, bench.SESS_INCREASE_GAUSSIAN, b)
	}
}

func benchmarkLoad(param int, rl zkbanw.RevocationList, name string, b usefulbench.Benchmarker) {
	parent := fmt.Sprintf("%s--%v", name, param)

	b.Run("load-"+parent, func(b usefulbench.Benchmarker) {
		for b.Loop() {
			test.PrepareUpdateKeyCached(rl.Sizes(), parent)
		}
	})
}
