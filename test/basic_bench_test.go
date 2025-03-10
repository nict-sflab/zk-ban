package zkbantest

import (
	"testing"

	zkban "github.com/akakou/zk-ban"
	"github.com/akakou/zk-ban/circuit"
	"github.com/akakou/zk-ban/snark"
	zkbanw "github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
)

func BenchmarkAll(t *testing.B) {
	params := prepareParams()

	rl1 := zkbanw.EmptyConstantRevocationAddList(270, 130)
	rl2 := zkbanw.EmptyLinerRevocationAddList(270, 130*270)
	rl3 := zkbanw.EmptyConstantRevocationAddList(1, 130*270)

	joinCircuit, signCircuit, updateCircuit1 := prepareCircuit(rl1, false)
	_, _, updateCircuit2 := prepareCircuit(rl2, true)
	_, _, updateCircuit3 := prepareCircuit(rl3, true)

	var proof groth16.Proof
	var err error

	var joinAssign *circuit.JoinRequestCircuit

	t.Run("join req", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			proof, joinAssign, err = zkban.JoinRequest(params.period, joinCircuit.Prover())
			panicIfErr(err)
		}
	})

	t.Run("verify join req", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			wit, err := frontend.NewWitness(joinAssign, snark.EcCurve.ScalarField())
			panicIfErr(err)

			pubWit, err := wit.Public()
			panicIfErr(err)

			err = groth16.Verify(proof, joinCircuit.VerifyKey, pubWit)
			panicIfErr(err)
		}
	})

	var signAssign *circuit.SignCircuit
	t.Run("sign", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			proof, signAssign, err = zkban.Sign(params.m, params.cnt, params.signer(), signCircuit.Prover())
			panicIfErr(err)
		}
	})

	t.Run("verify", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			wit, err := frontend.NewWitness(signAssign, snark.EcCurve.ScalarField())
			panicIfErr(err)

			pubWit, err := wit.Public()
			panicIfErr(err)

			err = groth16.Verify(proof, signCircuit.VerifyKey, pubWit)
			panicIfErr(err)
		}
	})

	var updateAssign *circuit.UpdateCircuit
	t.Run("update-req (constant)", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, proof, updateAssign, err = zkban.UpdateRequest(params.nextPeriod, params.signer(), rl1, params.gpk, updateCircuit1.Prover())
			panicIfErr(err)
		}
	})

	t.Run("update-verify (constant)", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			wit, err := frontend.NewWitness(updateAssign, snark.EcCurve.ScalarField())
			panicIfErr(err)

			pubWit, err := wit.Public()
			panicIfErr(err)

			params.gsk.IssueCredential(params.upk)

			err = groth16.Verify(proof, updateCircuit1.VerifyKey, pubWit)
			panicIfErr(err)
		}
	})

	t.Run("update-req (linear)", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, proof, updateAssign, err = zkban.UpdateRequest(params.nextPeriod, params.signer(), rl2, params.gpk, updateCircuit2.Prover())
			panicIfErr(err)
		}
	})

	t.Run("update-verify (linear)", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			wit, err := frontend.NewWitness(updateAssign, snark.EcCurve.ScalarField())
			panicIfErr(err)

			pubWit, err := wit.Public()
			panicIfErr(err)

			params.gsk.IssueCredential(params.upk)

			err = groth16.Verify(proof, updateCircuit2.VerifyKey, pubWit)
			panicIfErr(err)
		}
	})

	t.Run("update-req (one session)", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, proof, updateAssign, err = zkban.UpdateRequest(params.nextPeriod, params.signer(), rl3, params.gpk, updateCircuit3.Prover())
			panicIfErr(err)
		}
	})

	t.Run("update-verify (one session)", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			wit, err := frontend.NewWitness(updateAssign, snark.EcCurve.ScalarField())
			panicIfErr(err)

			pubWit, err := wit.Public()
			panicIfErr(err)

			params.gsk.IssueCredential(params.upk)

			err = groth16.Verify(proof, updateCircuit3.VerifyKey, pubWit)
			panicIfErr(err)
		}
	})
}
