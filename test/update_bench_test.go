package zkbantest

import (
	"fmt"
	"testing"

	zkban "github.com/akakou/zk-ban"
	"github.com/akakou/zk-ban/precomputes"
	zkbanw "github.com/akakou/zk-ban/witness"
	bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381"
	"github.com/consensys/gnark/backend/witness"
	// 	groth16_bls12381 "github.com/consensys/gnark/backend/groth16/bls12-381"
	// fr_bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
)

func BenchmarkUpdate(t *testing.B) {
	max := 10
	base := 600
	session := 180

	for i := 1; i <= max; i++ {
		benchmarkUpdate(session, base*i, t)
	}
}

func benchmarkUpdate(a, b int, t *testing.B) {
	params := prepareParams()

	var pubWit witness.Witness
	var err error

	var rl zkbanw.RevocationList
	var tag = ""
	rl = EmptyUniformRevocationList(a, b)
	tag = fmt.Sprintf("%d,%d,%d,%v", a*b, a, b, "constant")

	_, _, updateCircuit := prepareCircuit(rl, true)
	tag += fmt.Sprintf("-%v", updateCircuit.ConstraintSystem.GetNbPublicVariables())

	var proof *zkban.UpdateRequest
	t.Run("Prove ,"+tag, func(b *testing.B) {
		for b.Loop() {
			proof, err = zkban.RequestUpdate(params.nextPeriod, params.signer(), rl, params.gpk, updateCircuit.Prover())
			panicIfErr(err)
		}
	})

	vk, err := precomputes.NewUpdateVerificationKeyBLS12381(updateCircuit.VerifyKey)
	panicIfErr(err)

	var prepare *bls12381.G1Jac
	t.Run("Verify-Precomputes,"+tag, func(b *testing.B) {
		for b.Loop() {
			prepare, err = vk.PreparePublicInputs(pubWit)
			panicIfErr(err)
		}
	})

	t.Run("Verify ,"+tag, func(b *testing.B) {
		for b.Loop() {
			err = vk.VerifyPrepared(prepare, proof, params.nextPeriod, params.period)
			panicIfErr(err)

		}
	})
}
