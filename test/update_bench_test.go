package zkbantest

import (
	"fmt"
	"testing"

	zkban "github.com/akakou/zk-ban"
	"github.com/akakou/zk-ban/snark"
	zkbanw "github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/frontend"
	// 	groth16_bls12381 "github.com/consensys/gnark/backend/groth16/bls12-381"
	// fr_bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
)

func BenchmarkUpdate(t *testing.B) {
	max := 10
	base := 6000
	session := 270

	for i := 1; i <= max; i++ {
		benchmarkUpdate(session, base*i, t)
	}
}

func benchmarkUpdate(a, b int, t *testing.B) {
	params := prepareParams()

	var proof groth16.Proof
	var pubWit witness.Witness
	var err error

	var rl zkbanw.RevocationList
	var tag = ""
	rl = zkbanw.EmptyConstantRevocationAddList(a, b)
	tag = fmt.Sprintf("%d,%d,%d,%v", a*b, a, b, "constant")

	_, _, updateCircuit := prepareCircuit(rl, true)
	tag += fmt.Sprintf("-%v", updateCircuit.ConstraintSystem.GetNbPublicVariables())

	t.Run("Prove ,"+tag, func(b *testing.B) {
		for range b.N {
			_, proof2, updateAssign, err := zkban.UpdateRequest(params.nextPeriod, params.signer(), rl, params.gpk, updateCircuit.Prover())
			panicIfErr(err)

			proof = proof2
			wit, err := frontend.NewWitness(updateAssign, snark.EcCurve.ScalarField())
			panicIfErr(err)

			pubWit, err = wit.Public()
			panicIfErr(err)
		}
	})

	// inputs, err := groth16_bls12381.PreparePublicInputs(updateCircuit.VerifyKey.(*groth16_bls12381.VerifyingKey), pubWit.Vector().(fr_bls12381.Vector))
	// panicIfErr(err)

	t.Run("Verify ,"+tag, func(b *testing.B) {
		for range b.N {
			params.gsk.IssueCredential(params.upk)

			// err = groth16_bls12381.VerifyPrepared(proof.(*groth16_bls12381.Proof), updateCircuit.VerifyKey.(*groth16_bls12381.VerifyingKey), inputs)
			// panicIfErr(err)

			err = groth16.Verify(proof, updateCircuit.VerifyKey, pubWit)
			panicIfErr(err)
		}
	})
}
