package latency

import (
	"fmt"
	"runtime"

	zkban "github.com/akakou/zk-ban"
	"github.com/akakou/zk-ban/bench"
	"github.com/akakou/zk-ban/bench/latency/utils/usefulbench"
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
var max = 10

var baseSessionNum = 300

// var baseNymNum = 108_000 * 50
var baseNymNum = 30000

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

	if OnlyUniform {
		return
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
	if OnlyUniform {
		return
	}
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
