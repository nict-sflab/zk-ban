package zkbantest

import (
	"testing"

	gnarkprecomputes "github.com/akakou/gnark-precomputes"
	zkban "github.com/akakou/zk-ban"
	"github.com/akakou/zk-ban/circuit"
	zkbanc "github.com/akakou/zk-ban/circuit"
	bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381"
	"github.com/consensys/gnark/backend/witness"
)

func BenchmarkAll(t *testing.B) {
	params := prepareParams()

	rl1 := EmptyConstantRevocationAddList(270, 130)

	joinCircuit, signCircuit, updateCircuit1 := prepareCircuit(rl1, false)

	var err error

	var joinReq *zkban.JoinRequest
	t.Run("join req", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			joinReq, _, err = zkban.RequestJoin(params.period, joinCircuit.Prover())
			panicIfErr(err)
		}
	})

	t.Run("verify join req", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			err := joinReq.Verify(params.period, joinCircuit.VerifyKey)
			panicIfErr(err)
		}
	})

	var signature *zkban.Signature
	t.Run("sign", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			signature, err = zkban.Sign(params.m, params.cnt, params.signer(), params.gpk, signCircuit.Prover())
			panicIfErr(err)
		}
	})

	t.Run("verify", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			signature.Verify(params.m, params.cnt, params.period, params.gpk, signCircuit.VerifyKey)
		}
	})

	var updateRequest *zkban.UpdateRequest
	t.Run("update-req (constant)", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			updateRequest, err = zkban.RequestUpdate(params.nextPeriod, params.signer(), rl1, params.gpk, updateCircuit1.Prover())
			panicIfErr(err)
		}
	})

	vk, err := gnarkprecomputes.FromBLS12381GnarkKey(updateCircuit1.VerifyKey, updateCircuit1.Circuit.(*zkbanc.UpdateCircuit))
	panicIfErr(err)

	var prepare *bls12381.G1Jac
	var pubWit witness.Witness
	t.Run("update-verify-precomputes", func(b *testing.B) {
		for range b.N {
			pubWit, err = circuit.NewPublicUpdateCircuitWitness(params.nextPeriod, updateRequest.PublicKey, updateRequest.UpdateTicket, params.period, rl1, params.gpk)
			panicIfErr(err)

			prepare, err = vk.PreparePublicInputs(pubWit)
			panicIfErr(err)
		}
	})

	t.Run("update-verify (constant)", func(b *testing.B) {
		b.ResetTimer()

		for range b.N {
			params.gsk.IssueCredential(params.upk)
			err = vk.VerifyPrepared(updateRequest.Proof, pubWit, prepare)
			panicIfErr(err)
		}
	})

}
