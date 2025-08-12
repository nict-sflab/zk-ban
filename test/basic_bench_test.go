package zkbantest

import (
	"testing"

	zkban "github.com/akakou/zk-ban"
	"github.com/akakou/zk-ban/precomputes"
	bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381"
)

func BenchmarkAll(t *testing.B) {
	params := prepareParams()

	rl1 := EmptyUniformRevocationAddList(60, 300)

	joinCircuit, signCircuit, updateCircuit1 := prepareCircuit(rl1, false)

	var err error

	var joinReq *zkban.JoinRequest
	t.Run("join req", func(b *testing.B) {
		for b.Loop() {
			joinReq, _, err = zkban.RequestJoin(params.period, joinCircuit.Prover())
			panicIfErr(err)
		}
	})

	t.Run("verify join req", func(b *testing.B) {
		for b.Loop() {
			err := joinReq.Verify(params.period, joinCircuit.VerifyKey)
			panicIfErr(err)
		}
	})

	t.Run("issue credential", func(b *testing.B) {
		for b.Loop() {
			params.gsk.IssueCredential(params.upk)
		}
	})

	var signature *zkban.Signature
	t.Run("sign", func(b *testing.B) {
		for b.Loop() {
			signature, err = zkban.Sign(params.m, params.cnt, params.signer(), params.gpk, signCircuit.Prover())
			panicIfErr(err)
		}
	})

	t.Run("verify", func(b *testing.B) {
		for b.Loop() {
			signature.Verify(params.m, params.cnt, params.period, params.gpk, signCircuit.VerifyKey)
		}
	})

	var updateRequest *zkban.UpdateRequest
	t.Run("update-req (constant)", func(b *testing.B) {
		for b.Loop() {
			updateRequest, err = zkban.RequestUpdate(params.nextPeriod, params.signer(), rl1, params.gpk, updateCircuit1.Prover())
			panicIfErr(err)
		}
	})

	vk, err := precomputes.NewUpdateVerificationKeyBLS12381(updateCircuit1.VerifyKey)
	panicIfErr(err)

	var prepared **bls12381.G1Jac
	t.Run("update-verify-precomputes", func(b *testing.B) {
		for b.Loop() {
			prepared, err = vk.PrecomputeVerify(rl1, params.gpk)
			panicIfErr(err)
		}
	})

	t.Run("update-verify (constant)", func(b *testing.B) {
		for b.Loop() {
			err = vk.VerifyPrepared(*prepared, updateRequest, params.nextPeriod, params.period)
			panicIfErr(err)
		}
	})
}
