package latency

import (
	"runtime"

	"github.com/akakou/zk-ban/bench"
	"github.com/akakou/zk-ban/bench/latency/utils/usefulbench"
	"github.com/akakou/zk-ban/test"
	// 	groth16_bls12381 "github.com/consensys/gnark/backend/groth16/bls12-381"
	// fr_bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
)

func BenchmarkForCompare(b usefulbench.Benchmarker) {
	max := 10

	baseSessionNum := 60
	baseNymNum := 108_00

	ub, ok := b.(*usefulbench.UsefulBenchmaker)
	alpha := 1
	beta := 25

	if ok {
		ub.Result["env"]["baseSessionNum"] = int64(baseSessionNum)
		ub.Result["env"]["baseNymNum"] = int64(baseNymNum)
		ub.Result["env"]["alpha"] = int64(alpha)
		ub.Result["env"]["beta"] = int64(beta)
	}

	// increase nym
	for i := 1; i <= max; i++ {
		runtime.GC()

		nymNum := baseNymNum * i * alpha
		rl := test.EmptyUniformRevocationList(baseSessionNum, nymNum)
		benchmarkUpdate(nymNum, rl, bench.NYM_INCREASE_UNIFORM, b)
	}

	// for i := 1; i <= max; i++ {
	// 	runtime.GC()

	// 	nymNum := baseNymNum * i * alpha
	// 	rl := test.EmptyProportionalRevocationList(baseSessionNum, nymNum)
	// 	benchmarkForCompare(nymNum, rl, bench.NYM_INCREASE_PROPORTIONAL, b)
	// }

	// for i := 1; i <= max; i++ {
	// 	runtime.GC()

	// 	nymNum := baseNymNum * i * alpha
	// 	rl := test.EmptyGaussianRevocationList(baseSessionNum, nymNum)
	// 	benchmarkForCompare(nymNum, rl, bench.NYM_INCREASE_GAUSSIAN, b)
	// }
}

// func benchmarkForCompare(param int, rl zkbanw.RevocationList, name string, b usefulbench.Benchmarker) {
// 	params := test.PrepareParams()

// 	var err error

// 	parent := fmt.Sprintf("%s:%v", name, param)
// 	prover, veirifer := test.PrepareUpdateKeyCached(rl.Sizes(), parent)

// 	var proof *zkban.UpdateRequest
// 	b.Run("prove-"+parent, func(b usefulbench.Benchmarker) {
// 		runtime.GOMAXPROCS(1)
// 		for b.Loop() {
// 			proof, err = zkban.RequestUpdate(params.NextPeriod, params.Signer(), rl, params.GPK, prover)
// 			test.PanicIfErr(err)
// 		}
// 		runtime.GOMAXPROCS(runtime.NumCPU())
// 	})

// 	vk, err := precomputes.NewUpdateVerificationKeyBLS12381(veirifer.VerifyKey.VerifyingKey)
// 	test.PanicIfErr(err)

// 	var prepare **bls12381.G1Jac
// 	b.Run("precomputes-"+parent, func(b usefulbench.Benchmarker) {
// 		runtime.GOMAXPROCS(1)

// 		for b.Loop() {
// 			prepare, err = vk.PrecomputeVerify(rl, params.GPK)
// 			test.PanicIfErr(err)
// 		}
// 	})

// 	b.Run("verify-"+parent, func(b usefulbench.Benchmarker) {
// 		for b.Loop() {
// 			err = vk.VerifyPrepared(*prepare, proof, params.NextPeriod, params.Period)
// 			test.PanicIfErr(err)
// 		}
// 	})
// }
