package test

import (
	zkban "github.com/akakou/zk-ban"
	"github.com/akakou/zk-ban/precomputes"
	"github.com/akakou/zk-ban/test/utils/usefulbench"
	"github.com/akakou/zk-ban/witness"
	bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381"
)

func BenchmarkBaseline(b usefulbench.Benchmarker) {
	params := PrepareParams()

	baseNym := 108_000
	baseSess := 60

	ub, ok := b.(*usefulbench.UsefulBenchmaker)

	if ok {
		ub.Result["env"]["baseSessionNum"] = int64(baseSess)
		ub.Result["env"]["baseNymNum"] = int64(baseNym)
	}

	joinCircuit, signCircuit, _ := PrepareCircuit(witness.EmptyRevocationList([]int{}), false)
	var err error

	var joinReq *zkban.JoinRequest
	b.Run("baseline:request-join", func(b usefulbench.Benchmarker) {
		for b.Loop() {
			joinReq, _, err = zkban.RequestJoin(params.Period, joinCircuit.Prover())
			PanicIfErr(err)
		}
	})

	b.Run("baseline:verify-joinreq", func(b usefulbench.Benchmarker) {
		for b.Loop() {
			err := joinReq.Verify(params.Period, joinCircuit.VerifyKey)
			PanicIfErr(err)
		}
	})

	b.Run("baseline:issue credential", func(b usefulbench.Benchmarker) {
		for b.Loop() {
			params.GSK.IssueCredential(params.UPK)
		}
	})

	var signature *zkban.Signature
	b.Run("baseline:sign", func(b usefulbench.Benchmarker) {
		for b.Loop() {
			signature, err = zkban.Sign(params.M, params.CNT, params.Signer(), params.GPK, signCircuit.Prover())
			PanicIfErr(err)
		}
	})

	b.Run("baseline:verify", func(b usefulbench.Benchmarker) {
		for b.Loop() {
			signature.Verify(params.M, params.CNT, params.Period, params.GPK, signCircuit.VerifyKey)
		}
	})

	rlU := EmptyUniformRevocationList(baseSess, baseNym)
	benchmarkBasicUpdate(rlU, "uniform", &params, b)

	rlP := EmptyProportionalRevocationList(baseSess, baseNym)
	benchmarkBasicUpdate(rlP, "proportional", &params, b)

	rlG := EmptyGaussianRevocationList(baseSess, baseNym)
	benchmarkBasicUpdate(rlG, "gaussian", &params, b)

}

func benchmarkBasicUpdate(rl witness.RevocationList, name string, params *TestParams, b usefulbench.Benchmarker) {
	_, _, updateCircuit := PrepareCircuit(rl, false)

	var err error
	var updateRequest *zkban.UpdateRequest
	b.Run("baseline:update-req-"+name, func(b usefulbench.Benchmarker) {
		for b.Loop() {
			updateRequest, err = zkban.RequestUpdate(params.NextPeriod, params.Signer(), rl, params.GPK, updateCircuit.Prover())
			PanicIfErr(err)
		}
	})

	vk, err := precomputes.NewUpdateVerificationKeyBLS12381(updateCircuit.VerifyKey)
	PanicIfErr(err)

	var prepared **bls12381.G1Jac
	b.Run("baseline:update-precomputes-"+name, func(b usefulbench.Benchmarker) {
		for b.Loop() {
			prepared, err = vk.PrecomputeVerify(rl, params.GPK)
			PanicIfErr(err)
		}
	})

	b.Run("baseline:update-verify-"+name, func(b usefulbench.Benchmarker) {
		for b.Loop() {
			err = vk.VerifyPrepared(*prepared, updateRequest, params.NextPeriod, params.Period)
			PanicIfErr(err)
		}
	})
}
