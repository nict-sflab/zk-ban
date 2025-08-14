package test

import (
	"time"

	zkban "github.com/akakou/zk-ban"
	"github.com/akakou/zk-ban/precomputes"
	"github.com/akakou/zk-ban/test/utils/usefulbench"
	"github.com/akakou/zk-ban/witness"
	bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381"
)

func BenchmarkBaseline(b usefulbench.Benchmarker) {
	params := prepareParams()

	joinCircuit, signCircuit, _ := prepareCircuit(witness.EmptyRevocationList([]int{}), false)
	var err error

	var joinReq *zkban.JoinRequest
	b.Run("baseline:request-join", func(b usefulbench.Benchmarker) {
		for b.Loop() {
			joinReq, _, err = zkban.RequestJoin(params.period, joinCircuit.Prover())
			panicIfErr(err)
		}
	})

	b.Run("baseline:verify-joinreq", func(b usefulbench.Benchmarker) {
		for b.Loop() {
			err := joinReq.Verify(params.period, joinCircuit.VerifyKey)
			panicIfErr(err)
		}
	})

	b.Run("baseline:issue credential", func(b usefulbench.Benchmarker) {
		for b.Loop() {
			params.gsk.IssueCredential(params.upk)
		}
	})

	var signature *zkban.Signature
	b.Run("baseline:sign", func(b usefulbench.Benchmarker) {
		for b.Loop() {
			signature, err = zkban.Sign(params.m, params.cnt, params.signer(), params.gpk, signCircuit.Prover())
			panicIfErr(err)
		}
	})

	b.Run("baseline:verify", func(b usefulbench.Benchmarker) {
		for b.Loop() {
			signature.Verify(params.m, params.cnt, params.period, params.gpk, signCircuit.VerifyKey)
		}
	})

	baseNym := 108000
	baseSess := 180

	rlU := EmptyUniformRevocationList(baseSess, baseNym)
	benchmarkBasicUpdate(rlU, "uniform", &params, b)
	time.Sleep(time.Second * 3)

	rlP := EmptyProportionalRevocationList(baseSess, baseNym)
	benchmarkBasicUpdate(rlP, "proportional", &params, b)
	time.Sleep(time.Second * 3)

	rlG := EmptyGaussianRevocationList(baseSess, baseNym)
	benchmarkBasicUpdate(rlG, "gaussian", &params, b)

}

func benchmarkBasicUpdate(rl witness.RevocationList, name string, params *TestParams, b usefulbench.Benchmarker) {
	_, _, updateCircuit := prepareCircuit(rl, false)

	var err error
	var updateRequest *zkban.UpdateRequest
	b.Run("baseline:update-req-"+name, func(b usefulbench.Benchmarker) {
		for b.Loop() {
			updateRequest, err = zkban.RequestUpdate(params.nextPeriod, params.signer(), rl, params.gpk, updateCircuit.Prover())
			panicIfErr(err)
		}
	})

	vk, err := precomputes.NewUpdateVerificationKeyBLS12381(updateCircuit.VerifyKey)
	panicIfErr(err)

	var prepared **bls12381.G1Jac
	b.Run("baseline:update-precomputes-"+name, func(b usefulbench.Benchmarker) {
		for b.Loop() {
			prepared, err = vk.PrecomputeVerify(rl, params.gpk)
			panicIfErr(err)
		}
	})

	b.Run("baseline:update-verify-"+name, func(b usefulbench.Benchmarker) {
		for b.Loop() {
			err = vk.VerifyPrepared(*prepared, updateRequest, params.nextPeriod, params.period)
			panicIfErr(err)
		}
	})
}
