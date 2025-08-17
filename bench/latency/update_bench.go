package latency

import (
	"fmt"

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

func BenchmarkUpdate(b usefulbench.Benchmarker) {
	max := 10

	baseSessionNum := 60
	baseNymNum := 108_000

	ub, ok := b.(*usefulbench.UsefulBenchmaker)
	alpha := 1
	beta := 50

	if ok {
		ub.Result["env"]["baseSessionNum"] = int64(baseSessionNum)
		ub.Result["env"]["baseNymNum"] = int64(baseNymNum)
		ub.Result["env"]["alpha"] = int64(alpha)
		ub.Result["env"]["beta"] = int64(beta)
	}

	// increase nym
	for i := 1; i <= max; i++ {
		nymNum := baseNymNum * i * alpha
		rl := test.EmptyUniformRevocationList(baseSessionNum, nymNum)
		benchmarkUpdate(nymNum, rl, bench.NYM_INCREASE_UNIFORM, b)
	}

	for i := 1; i <= max; i++ {
		nymNum := baseNymNum * i * alpha
		rl := test.EmptyProportionalRevocationList(baseSessionNum, nymNum)
		benchmarkUpdate(nymNum, rl, bench.NYM_INCREASE_PROPORTIONAL, b)
	}

	for i := 1; i <= max; i++ {
		nymNum := baseNymNum * i * alpha
		rl := test.EmptyGaussianRevocationList(baseSessionNum, nymNum)
		benchmarkUpdate(nymNum, rl, bench.NYM_INCREASE_GAUSSIAN, b)
	}

	// increase sessionNumber
	for i := 1; i <= max; i++ {
		sessionNum := baseSessionNum * i * beta
		rl := test.EmptyUniformRevocationList(sessionNum, baseNymNum)
		benchmarkUpdate(sessionNum, rl, bench.SESS_INCREASE_UNIFORM, b)
	}

	for i := 1; i <= max; i++ {
		sessionNum := baseSessionNum * i * beta
		rl := test.EmptyProportionalRevocationList(sessionNum, baseNymNum)
		benchmarkUpdate(sessionNum, rl, bench.SESS_INCREASE_PROPORTIONAL, b)
	}

	for i := 1; i <= max; i++ {
		sessionNum := baseSessionNum * i * beta
		rl := test.EmptyGaussianRevocationList(sessionNum, baseNymNum)
		benchmarkUpdate(sessionNum, rl, bench.SESS_INCREASE_GAUSSIAN, b)
	}
}

func benchmarkUpdate(param int, rl zkbanw.RevocationList, name string, b usefulbench.Benchmarker) {
	params := test.PrepareParams()

	var err error

	parent := fmt.Sprintf("%s:%v", name, param)
	prover, veirifer := test.PrepareUpdateKeyCached(rl.Sizes(), parent)

	var proof *zkban.UpdateRequest
	b.Run("prove-"+parent, func(b usefulbench.Benchmarker) {
		for b.Loop() {
			proof, err = zkban.RequestUpdate(params.NextPeriod, params.Signer(), rl, params.GPK, prover)
			test.PanicIfErr(err)
		}
	})

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
