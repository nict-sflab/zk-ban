package latency

import (
	"fmt"
	"runtime"

	zkban "github.com/akakou/zk-ban"
	"github.com/akakou/zk-ban/bench"
	"github.com/akakou/zk-ban/bench/latency/utils/usefulbench"
	"github.com/akakou/zk-ban/precomputes"
	"github.com/akakou/zk-ban/test"
	zkbanw "github.com/akakou/zk-ban/witness"
	bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381"
	// 	groth16_bls12381 "github.com/consensys/gnark/backend/groth16/bls12-381"
	// fr_bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
)

var alpha = 1
var beta = 25
var max = 10

var baseSessionNum = 60
var baseNymNum = 108_000

func logIfUsefulBenchInScaleBench(b usefulbench.Benchmarker) {
	ub, ok := b.(*usefulbench.UsefulBenchmaker)

	if ok {
		ub.Result["env"]["baseSessionNum"] = int64(baseSessionNum)
		ub.Result["env"]["baseNymNum"] = int64(baseNymNum)
		ub.Result["env"]["alpha"] = int64(alpha)
		ub.Result["env"]["beta"] = int64(beta)
		ub.Result["env"]["max"] = int64(max)
	}
}

func BenchmarkScalability(b usefulbench.Benchmarker) {
	BenchmarkNymScalability(b)
	BenchSessScalability(b)
}

func BenchmarkNymScalability(b usefulbench.Benchmarker) {
	logIfUsefulBenchInScaleBench(b)

	// increase nym
	for i := 1; i <= max; i++ {
		runtime.GC()

		nymNum := baseNymNum * i * alpha
		rl := test.EmptyUniformRevocationList(baseSessionNum, nymNum)
		benchmarkUpdate(nymNum, rl, bench.NYM_INCREASE_UNIFORM, b)
	}

	for i := 1; i <= max; i++ {
		runtime.GC()

		nymNum := baseNymNum * i * alpha
		rl := test.EmptyProportionalRevocationList(baseSessionNum, nymNum)
		benchmarkUpdate(nymNum, rl, bench.NYM_INCREASE_PROPORTIONAL, b)
	}

	for i := 1; i <= max; i++ {
		runtime.GC()

		nymNum := baseNymNum * i * alpha
		rl := test.EmptyGaussianRevocationList(baseSessionNum, nymNum)
		benchmarkUpdate(nymNum, rl, bench.NYM_INCREASE_GAUSSIAN, b)
	}
}

func BenchSessScalability(b usefulbench.Benchmarker) {
	// increase sessionNumber
	BenchSessUniformScalability(b)
	BenchSessProportionalScalability(b)
	BenchSessGaussScalability(b)

}

func BenchSessUniformScalability(b usefulbench.Benchmarker) {
	logIfUsefulBenchInScaleBench(b)

	for i := 1; i <= max; i++ {
		runtime.GC()

		sessionNum := baseSessionNum * i * beta
		rl := test.EmptyUniformRevocationList(sessionNum, baseNymNum)
		benchmarkUpdate(sessionNum, rl, bench.SESS_INCREASE_UNIFORM, b)
	}

}

func BenchSessProportionalScalability(b usefulbench.Benchmarker) {
	logIfUsefulBenchInScaleBench(b)

	for i := 1; i <= max; i++ {
		runtime.GC()

		sessionNum := baseSessionNum * i * beta
		rl := test.EmptyProportionalRevocationList(sessionNum, baseNymNum)
		benchmarkUpdate(sessionNum, rl, bench.SESS_INCREASE_PROPORTIONAL, b)
	}
}

func BenchSessGaussScalability(b usefulbench.Benchmarker) {
	logIfUsefulBenchInScaleBench(b)

	for i := 1; i <= max; i++ {
		runtime.GC()

		sessionNum := baseSessionNum * i * beta
		rl := test.EmptyGaussianRevocationList(sessionNum, baseNymNum)
		benchmarkUpdate(sessionNum, rl, bench.SESS_INCREASE_GAUSSIAN, b)
	}
}

var NoParallel = false
var SkipVerify = false

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

	var proof *zkban.UpdateRequest
	b.Run("prove-"+parent, func(b usefulbench.Benchmarker) {
		for b.Loop() {
			proof, err = zkban.RequestUpdate(params.NextPeriod, params.Signer(), rl, params.GPK, prover)
			test.PanicIfErr(err)
		}
	})

	if SkipVerify {
		return
	}

	vk, err := precomputes.NewUpdateVerificationKeyBLS12381(veirifer.VerifyKey.VerifyingKey)
	test.PanicIfErr(err)

	var prepare **bls12381.G1Jac
	b.Run("precomputes-"+parent, func(b usefulbench.Benchmarker) {
		for b.Loop() {
			prepare, err = vk.PrecomputeVerify(rl, params.GPK)
			test.PanicIfErr(err)
		}
	})

	b.Run("verify-"+parent, func(b usefulbench.Benchmarker) {
		for b.Loop() {
			err = vk.VerifyPrepared(*prepare, proof, params.NextPeriod, params.Period)
			test.PanicIfErr(err)
		}
	})

}
