package zkbantest

import (
	"fmt"
	"testing"

	zkban "github.com/akakou/zk-ban"
	"github.com/akakou/zk-ban/precomputes"
	zkbanw "github.com/akakou/zk-ban/witness"
	bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381"
	// 	groth16_bls12381 "github.com/consensys/gnark/backend/groth16/bls12-381"
	// fr_bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
)

func BenchmarkUpdate(t *testing.B) {
	max := 10

	baseSessionNum := 180
	baseNymNum := baseSessionNum * 900

	a := 1

	// increase nym
	for i := 1; i <= max; i++ {
		nymNum := baseNymNum * i * a
		rl := EmptyUniformRevocationList(baseSessionNum, nymNum)
		benchmarkUpdate(baseSessionNum, nymNum, rl, "uniform", t)
	}

	for i := 1; i <= max; i++ {
		nymNum := baseNymNum * i * a
		rl := EmptyProportionalRevocationList(baseSessionNum, nymNum)
		benchmarkUpdate(baseSessionNum, nymNum, rl, "proportional", t)
	}

	for i := 1; i <= max; i++ {
		nymNum := baseNymNum * i * a
		rl := EmptyGaussianRevocationList(baseSessionNum, nymNum)
		benchmarkUpdate(baseSessionNum, nymNum, rl, "gaussian", t)
	}

	// increase sessionNumber
	b := 1
	for i := 1; i <= max; i++ {
		sessionNum := baseSessionNum * i * b
		rl := EmptyUniformRevocationList(sessionNum, baseNymNum)
		benchmarkUpdate(sessionNum, baseNymNum, rl, "uniform", t)
	}

	for i := 1; i <= max; i++ {
		sessionNum := baseSessionNum * i * b
		rl := EmptyProportionalRevocationList(sessionNum, baseNymNum)
		benchmarkUpdate(sessionNum, baseNymNum, rl, "proportional", t)
	}

	for i := 1; i <= max; i++ {
		sessionNum := baseSessionNum * i * b
		rl := EmptyGaussianRevocationList(sessionNum, baseNymNum)
		benchmarkUpdate(sessionNum, baseNymNum, rl, "gaussian", t)
	}
}

func benchmarkUpdate(sessionNumberSize, nymNum int, rl zkbanw.RevocationList, name string, t *testing.B) {
	params := prepareParams()

	var err error

	_, _, updateCircuit := prepareCircuit(rl, true)
	name += fmt.Sprintf("%s: %v-%v", name, sessionNumberSize, nymNum)

	var proof *zkban.UpdateRequest
	t.Run("Prove: "+name, func(b *testing.B) {
		for b.Loop() {
			proof, err = zkban.RequestUpdate(params.nextPeriod, params.signer(), rl, params.gpk, updateCircuit.Prover())
			panicIfErr(err)
		}
	})

	vk, err := precomputes.NewUpdateVerificationKeyBLS12381(updateCircuit.VerifyKey)
	panicIfErr(err)

	var prepare **bls12381.G1Jac
	t.Run("Verify-Precomputes: "+name, func(b *testing.B) {
		for b.Loop() {
			prepare, err = vk.PrecomputeVerify(rl, params.gpk)
			panicIfErr(err)
		}
	})

	t.Run("Verify: "+name, func(b *testing.B) {
		for b.Loop() {
			err = vk.VerifyPrepared(*prepare, proof, params.nextPeriod, params.period)
			panicIfErr(err)

		}
	})
}
