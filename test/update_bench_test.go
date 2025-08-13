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
	base := 600
	session := 180

	for i := 1; i <= max; i++ {
		nymSize := session * base * i
		rl := EmptyUniformRevocationList(session, nymSize)
		benchmarkUpdate(session, base*session, rl, "uniform", t)
	}

	for i := 1; i <= max; i++ {
		nymSize := session * base * i
		rl := EmptyProportionalRevocationList(session, nymSize)
		benchmarkUpdate(session, nymSize, rl, "proportional", t)
	}

	for i := 1; i <= max; i++ {
		nymSize := session * base * i
		rl := EmptyGaussianRevocationList(session, nymSize)
		benchmarkUpdate(session, nymSize, rl, "gaussian", t)
	}
}

func benchmarkUpdate(sessionSize, nymSize int, rl zkbanw.RevocationList, name string, t *testing.B) {
	params := prepareParams()

	var err error

	_, _, updateCircuit := prepareCircuit(rl, true)
	name += fmt.Sprintf("%s: %v-%v", name, sessionSize, nymSize)

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
