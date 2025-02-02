package zkban_test

import (
	"testing"

	zkban "github.com/akakou/zk-ban"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
)

func BenchmarkAll(t *testing.B) {
	params := prepare()

	var proof groth16.Proof
	var pubWit witness.Witness
	var err error

	t.Run("join req", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			proof, pubWit, _, _, err = zkban.JoinRequest(params.period, params.joinSnark.Prover())
			panicIfErr(err)
		}
	})

	t.Run("verify join req", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			err = groth16.Verify(proof, params.joinSnark.VerifyKey, pubWit)
			panicIfErr(err)
		}
	})

	t.Run("sign", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			proof, pubWit, err = zkban.Sign(params.m, params.bsn, params.signer(), params.gpk, params.signSnark.Prover())
			panicIfErr(err)
		}
	})

	t.Run("verify", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			err = groth16.Verify(proof, params.signSnark.VerifyKey, pubWit)
			panicIfErr(err)
		}
	})

	t.Run("update-req", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, proof, pubWit, err = zkban.Update(params.revocationList, params.nextPeriod, params.signer(), params.sessionName, params.gpk, params.updateSnark.Prover())
			panicIfErr(err)
		}
	})

	t.Run("update-verify", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			params.gsk.IssueCredential(params.upk)

			err = groth16.Verify(proof, params.updateSnark.VerifyKey, pubWit)
			panicIfErr(err)
		}
	})
}
