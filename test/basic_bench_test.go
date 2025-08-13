package zkbantest

import (
	"testing"

	zkban "github.com/akakou/zk-ban"
	"github.com/akakou/zk-ban/precomputes"
	"github.com/akakou/zk-ban/witness"
	bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381"
)

func BenchmarkAll(b *testing.B) {
	params := prepareParams()

	joinCircuit, signCircuit, _ := prepareCircuit(witness.EmptyRevocationList([]int{}), false)
	var err error

	var joinReq *zkban.JoinRequest
	b.Run("join req", func(b *testing.B) {
		for b.Loop() {
			joinReq, _, err = zkban.RequestJoin(params.period, joinCircuit.Prover())
			panicIfErr(err)
		}
	})

	b.Run("verify join req", func(b *testing.B) {
		for b.Loop() {
			err := joinReq.Verify(params.period, joinCircuit.VerifyKey)
			panicIfErr(err)
		}
	})

	b.Run("issue credential", func(b *testing.B) {
		for b.Loop() {
			params.gsk.IssueCredential(params.upk)
		}
	})

	var signature *zkban.Signature
	b.Run("sign", func(b *testing.B) {
		for b.Loop() {
			signature, err = zkban.Sign(params.m, params.cnt, params.signer(), params.gpk, signCircuit.Prover())
			panicIfErr(err)
		}
	})

	b.Run("verify", func(b *testing.B) {
		for b.Loop() {
			signature.Verify(params.m, params.cnt, params.period, params.gpk, signCircuit.VerifyKey)
		}
	})

	rlU := EmptyUniformRevocationList(180, 90000)
	benchmarkBasicUpdate(rlU, "uniform", &params, b)
	rlP := EmptyProportionalRevocationList(180, 90000)
	benchmarkBasicUpdate(rlP, "proportional", &params, b)
	rlG := EmptyProportionalRevocationList(180, 90000)
	benchmarkBasicUpdate(rlG, "gaussian", &params, b)

}

func benchmarkBasicUpdate(rl witness.RevocationList, name string, params *TestParams, b *testing.B) {
	_, _, updateCircuit := prepareCircuit(rl, false)

	var err error
	var updateRequest *zkban.UpdateRequest
	b.Run("update-req: "+name, func(b *testing.B) {
		for b.Loop() {
			updateRequest, err = zkban.RequestUpdate(params.nextPeriod, params.signer(), rl, params.gpk, updateCircuit.Prover())
			panicIfErr(err)
		}
	})

	vk, err := precomputes.NewUpdateVerificationKeyBLS12381(updateCircuit.VerifyKey)
	panicIfErr(err)

	var prepared **bls12381.G1Jac
	b.Run("update-verify-precomputes: "+name, func(b *testing.B) {
		for b.Loop() {
			prepared, err = vk.PrecomputeVerify(rl, params.gpk)
			panicIfErr(err)
		}
	})

	b.Run("update-verify : "+name, func(b *testing.B) {
		for b.Loop() {
			err = vk.VerifyPrepared(*prepared, updateRequest, params.nextPeriod, params.period)
			panicIfErr(err)
		}
	})
}
