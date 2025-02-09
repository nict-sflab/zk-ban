package zkbantest

import (
	"testing"

	zkban "github.com/akakou/zk-ban"
	"github.com/akakou/zk-ban/commit"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
)

func BenchmarkAll(t *testing.B) {
	params := prepareParams()

	rl := commit.EmptyConstantRevocationAddList(270, 130)
	joinCircuit, signCircuit, updateCircuit := prepareCircuit(rl, false)

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
			proof, pubWit, err = zkban.Sign(params.m, params.bsn, params.signer(), signCircuit.Prover())
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

	t.Run("update-req", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, proof, pubWit, err = zkban.Update(params.nextPeriod, params.gpk, params.signer(), rl, updateCircuit.Prover())
			panicIfErr(err)
		}
	})

	t.Run("update-verify", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			params.gsk.IssueCredential(params.upk)

			err = groth16.Verify(proof, updateCircuit.VerifyKey, pubWit)
			panicIfErr(err)
		}
	})
}
