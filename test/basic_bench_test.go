package zkbantest

import (
	"testing"

	zkban "github.com/akakou/zk-ban"
	zkbanw "github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
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
	var pubWit witness.Witness
	var err error

	t.Run("join req", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			proof, pubWit, _, _, err = zkban.JoinRequest(params.period, joinCircuit.Prover())
			panicIfErr(err)
		}
	})

	t.Run("verify join req", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			err = groth16.Verify(proof, joinCircuit.VerifyKey, pubWit)
			panicIfErr(err)
		}
	})

	t.Run("sign", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			proof, pubWit, err = zkban.Sign(params.m, params.cnt, params.signer(), signCircuit.Prover())
			panicIfErr(err)
		}
	})

	t.Run("verify", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			err = groth16.Verify(proof, signCircuit.VerifyKey, pubWit)
			panicIfErr(err)
		}
	})

	t.Run("update-req (constant)", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, proof, pubWit, err = zkban.Update(params.nextPeriod, params.gpk, params.signer(), rl1, updateCircuit1.Prover())
			panicIfErr(err)
		}
	})

	t.Run("update-verify (constant)", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			params.gsk.IssueCredential(params.upk)

			err = groth16.Verify(proof, updateCircuit1.VerifyKey, pubWit)
			panicIfErr(err)
		}
	})

	t.Run("update-req (linear)", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, proof, pubWit, err = zkban.Update(params.nextPeriod, params.gpk, params.signer(), rl2, updateCircuit2.Prover())
			panicIfErr(err)
		}
	})

	t.Run("update-verify (linear)", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			params.gsk.IssueCredential(params.upk)

			err = groth16.Verify(proof, updateCircuit2.VerifyKey, pubWit)
			panicIfErr(err)
		}
	})

	t.Run("update-req (one session)", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, proof, pubWit, err = zkban.Update(params.nextPeriod, params.gpk, params.signer(), rl3, updateCircuit3.Prover())
			panicIfErr(err)
		}
	})

	t.Run("update-verify (one session)", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			params.gsk.IssueCredential(params.upk)

			err = groth16.Verify(proof, updateCircuit3.VerifyKey, pubWit)
			panicIfErr(err)
		}
	})
}
