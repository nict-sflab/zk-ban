package test

import (
	"fmt"

	zkban "github.com/akakou/zk-ban"
	"github.com/akakou/zk-ban/precomputes"
	"github.com/akakou/zk-ban/test/utils/usefulbench"
	zkbanw "github.com/akakou/zk-ban/witness"
	bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381"
	// 	groth16_bls12381 "github.com/consensys/gnark/backend/groth16/bls12-381"
	// fr_bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
)

func BenchmarkUpdate(b usefulbench.Benchmarker) {
	max := 5

	baseSessionNum := 180
	baseNymNum := baseSessionNum * 1080

	alpha := 1

	// increase nym
	for i := 1; i <= max; i++ {
		nymNum := baseNymNum * i * alpha
		rl := EmptyUniformRevocationList(baseSessionNum, nymNum)
		benchmarkUpdate(nymNum, rl, "nym-increase-uniform", b)
	}

	for i := 1; i <= max; i++ {
		nymNum := baseNymNum * i * alpha
		rl := EmptyProportionalRevocationList(baseSessionNum, nymNum)
		benchmarkUpdate(nymNum, rl, "nym-increase-proportional", b)
	}

	for i := 1; i <= max; i++ {
		nymNum := baseNymNum * i * alpha
		rl := EmptyGaussianRevocationList(baseSessionNum, nymNum)
		benchmarkUpdate(nymNum, rl, "nym-increase-gaussian", b)
	}

	// increase sessionNumber
	beta := 10
	for i := 1; i <= max; i++ {
		sessionNum := baseSessionNum * i * beta
		rl := EmptyUniformRevocationList(sessionNum, baseNymNum)
		benchmarkUpdate(sessionNum, rl, "sess-increase-uniform", b)
	}

	for i := 1; i <= max; i++ {
		sessionNum := baseSessionNum * i * beta
		rl := EmptyProportionalRevocationList(sessionNum, baseNymNum)
		benchmarkUpdate(sessionNum, rl, "sess-increase-proportional", b)
	}

	for i := 1; i <= max; i++ {
		sessionNum := baseSessionNum * i * beta
		rl := EmptyGaussianRevocationList(sessionNum, baseNymNum)
		benchmarkUpdate(sessionNum, rl, "sess-increase-gaussian", b)
	}
}

func benchmarkUpdate(param int, rl zkbanw.RevocationList, name string, b usefulbench.Benchmarker) {
	params := PrepareParams()

	var err error

	_, _, updateCircuit := PrepareCircuit(rl, true)
	name += fmt.Sprintf("%s:%v", name, param)

	var proof *zkban.UpdateRequest
	b.Run("prove-"+name, func(b usefulbench.Benchmarker) {
		for b.Loop() {
			proof, err = zkban.RequestUpdate(params.NextPeriod, params.Signer(), rl, params.GPK, updateCircuit.Prover())
			PanicIfErr(err)
		}
	})

	vk, err := precomputes.NewUpdateVerificationKeyBLS12381(updateCircuit.VerifyKey)
	PanicIfErr(err)

	var prepare **bls12381.G1Jac
	b.Run("precomputes-"+name, func(b usefulbench.Benchmarker) {
		for b.Loop() {
			prepare, err = vk.PrecomputeVerify(rl, params.GPK)
			PanicIfErr(err)
		}
	})

	b.Run("verify-"+name, func(b usefulbench.Benchmarker) {
		for b.Loop() {
			err = vk.VerifyPrepared(*prepare, proof, params.NextPeriod, params.Period)
			PanicIfErr(err)

		}
	})
}
