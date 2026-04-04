package latency

import (
	"fmt"
	"runtime"

	zkban "github.com/akakou/zk-ban"
	"github.com/akakou/zk-ban/bench"
	"github.com/akakou/zk-ban/bench/latency/utils/usefulbench"
	"github.com/akakou/zk-ban/circuit"
	"github.com/akakou/zk-ban/load"
	"github.com/akakou/zk-ban/precomputes"
	"github.com/akakou/zk-ban/test"
	zkbanw "github.com/akakou/zk-ban/witness"
	// 	groth16_bls12381 "github.com/consensys/gnark/backend/groth16/bls12-381"
	// fr_bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
)

var NoParallel = false
var SkipVerify = false
var OnlyUniform = false

var alpha = 1
var beta = 10
var gamma = 10
var max = 10

var baseSessNum = 5
var basePeriodNum = 300
var baseNymNum = 30000

func logIfUsefulBenchInScaleBench(b usefulbench.Benchmarker) {
	ub, ok := b.(*usefulbench.UsefulBenchmaker)

	if ok {
		ub.Result["env"]["basePeriodNum"] = int64(basePeriodNum)
		ub.Result["env"]["baseNymNum"] = int64(baseNymNum)
		ub.Result["env"]["alpha"] = int64(alpha)
		ub.Result["env"]["beta"] = int64(beta)
		ub.Result["env"]["max"] = int64(max)
	}
}

func BenchmarkScalability(b usefulbench.Benchmarker) {
	BenchmarkNymScalability(b)
	BenchPeriodScalability(b)
	BenchMaxSessScalability(b)
}

func BenchmarkNymScalability(b usefulbench.Benchmarker) {
	logIfUsefulBenchInScaleBench(b)

	// increase nym
	BenchmarkOneNymScalability(b, bench.NYM_INCREASE_UNIFORM, test.EmptyUniformRevocationList)
	if OnlyUniform {
		return
	}
	BenchmarkOneNymScalability(b, bench.NYM_INCREASE_PROPORTIONAL, test.EmptyProportionalRevocationList)
	BenchmarkOneNymScalability(b, bench.NYM_INCREASE_GAUSSIAN, test.EmptyGaussianRevocationList)

}
func BenchmarkOneNymScalability(b usefulbench.Benchmarker, name string, EmptyRevocationList func(int, int) zkbanw.RevocationList) {
	logIfUsefulBenchInScaleBench(b)

	// increase nym
	for i := 1; i <= max; i++ {
		runtime.GC()

		nymNum := baseNymNum * i * alpha
		rl := EmptyRevocationList(basePeriodNum, nymNum)
		benchmarkUpdate(nymNum, rl, bench.NYM_INCREASE_UNIFORM, b)
	}

}

func BenchPeriodScalability(b usefulbench.Benchmarker) {
	// increase periodNumber
	BenchOnePeriodScalability(b, bench.PERIOD_INCREASE_UNIFORM, test.EmptyUniformRevocationList)
	if OnlyUniform {
		return
	}
	BenchOnePeriodScalability(b, bench.PERIOD_INCREASE_PROPORTIONAL, test.EmptyProportionalRevocationList)
	BenchOnePeriodScalability(b, bench.PERIOD_INCREASE_GAUSSIAN, test.EmptyGaussianRevocationList)
}

func BenchOnePeriodScalability(b usefulbench.Benchmarker, name string, EmptyRevocationList func(int, int) zkbanw.RevocationList) {
	logIfUsefulBenchInScaleBench(b)

	for i := 1; i <= max; i++ {
		runtime.GC()

		periodNum := basePeriodNum * i * beta
		rl := EmptyRevocationList(periodNum, baseNymNum)
		benchmarkUpdate(periodNum, rl, bench.PERIOD_INCREASE_UNIFORM, b)
	}

}

func BenchMaxSessScalability(b usefulbench.Benchmarker) {
	// increase periodNumber
	BenchOneMaxSessScalability(b, bench.SESS_INCREASE_UNIFORM, test.EmptyUniformRevocationList)
	if OnlyUniform {
		return
	}
	BenchOneMaxSessScalability(b, bench.SESS_INCREASE_PROPORTIONAL, test.EmptyProportionalRevocationList)
	BenchOneMaxSessScalability(b, bench.SESS_INCREASE_GAUSSIAN, test.EmptyGaussianRevocationList)
}

func BenchOneMaxSessScalability(b usefulbench.Benchmarker, name string, EmptyRevocationList func(int, int) zkbanw.RevocationList) {
	logIfUsefulBenchInScaleBench(b)

	for i := 1; i <= max; i++ {
		runtime.GC()
		sessNum := baseSessNum * i * gamma

		circuit.MaxSession = sessNum

		rl := EmptyRevocationList(baseNymNum, baseNymNum)
		benchmarkUpdate(sessNum, rl, bench.PERIOD_INCREASE_PROPORTIONAL, b)
	}

}

func benchmarkUpdate(param int, rl zkbanw.RevocationList, name string, b usefulbench.Benchmarker) {
	params := test.PrepareParams()

	var err error

	parent := fmt.Sprintf("%s--%v", name, param)
	prover, veirifer := test.PrepareUpdateKeyCached(rl.Sizes(), parent)

	if NoParallel {
		runtime.GOMAXPROCS(1)
	}

	defer func() {
		if NoParallel {
			runtime.GOMAXPROCS(runtime.NumCPU())
		}
	}()

	b.Run("load-prove-"+parent, func(b usefulbench.Benchmarker) {
		for b.Loop() {
			_, err := load.LoadUserKey(parent, "update")
			test.PanicIfErr(err)
		}
	})

	b.Run("prove-"+parent, func(b usefulbench.Benchmarker) {
		for b.Loop() {
			_, err := zkban.RequestUpdate(params.NextPeriod, params.Signer(), rl, params.GPK, prover)
			test.PanicIfErr(err)
		}
	})
	proof, _ := zkban.RequestUpdate(params.NextPeriod, params.Signer(), rl, params.GPK, prover)
	test.PanicIfErr(err)

	if SkipVerify {
		return
	}

	b.Run("load-verify-"+parent, func(b usefulbench.Benchmarker) {
		for b.Loop() {
			_, err := load.LoadGroupManagerUpdateKey(parent)
			test.PanicIfErr(err)
		}
	})

	vk, err := precomputes.NewUpdateVerificationKeyBLS12381(*veirifer.VerifyKey)
	test.PanicIfErr(err)

	b.Run("precomputes-"+parent, func(b usefulbench.Benchmarker) {
		for b.Loop() {
			_, err := vk.PrecomputeVerify(rl, params.GPK)
			test.PanicIfErr(err)
		}
	})

	prepare, _ := vk.PrecomputeVerify(rl, params.GPK)

	b.Run("verify-"+parent, func(b usefulbench.Benchmarker) {
		for b.Loop() {
			err = vk.VerifyPrepared(*prepare, proof, params.NextPeriod, params.Period)
			test.PanicIfErr(err)
		}
	})

}
