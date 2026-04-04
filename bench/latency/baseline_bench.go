package latency

import (
	"runtime"

	zkban "github.com/akakou/zk-ban"
	"github.com/akakou/zk-ban/bench/latency/utils/usefulbench"
	"github.com/akakou/zk-ban/load"
	"github.com/akakou/zk-ban/precomputes"
	"github.com/akakou/zk-ban/test"
	"github.com/akakou/zk-ban/witness"
)

func BenchmarkBaseline(b usefulbench.Benchmarker) {
	params := test.PrepareParams()

	baseNym := 30000
	basePeriod := 30

	ub, ok := b.(*usefulbench.UsefulBenchmaker)

	if ok {
		ub.Result["env"]["basePeriodNum"] = int64(basePeriod)
		ub.Result["env"]["baseNymNum"] = int64(baseNym)
	}

	benchmarkBaseline(&params, b)

	rlU := test.EmptyUniformRevocationList(basePeriod, baseNym)
	benchmarkBasicUpdate(rlU, "uniform", &params, b)

	rlP := test.EmptyProportionalRevocationList(basePeriod, baseNym)
	benchmarkBasicUpdate(rlP, "proportional", &params, b)

	rlG := test.EmptyGaussianRevocationList(basePeriod, baseNym)
	benchmarkBasicUpdate(rlG, "gaussian", &params, b)
}

func benchmarkBaseline(params *test.TestParams, b usefulbench.Benchmarker) {
	joinCircuit, signCircuit, _ := test.PrepareCircuit(witness.EmptyRevocationList([]int{}))
	var err error

	if NoParallel {
		runtime.GOMAXPROCS(1)
	}

	b.Run("baseline--request-join", func(b usefulbench.Benchmarker) {
		for b.Loop() {
			_, _, err := zkban.RequestJoin(params.Period, joinCircuit.Prover())
			test.PanicIfErr(err)
		}
	})

	joinReq, _, _ := zkban.RequestJoin(params.Period, joinCircuit.Prover())

	b.Run("baseline--verify-joinreq", func(b usefulbench.Benchmarker) {
		for b.Loop() {
			err := joinReq.Verify(params.Period, joinCircuit.VerifyKey)
			test.PanicIfErr(err)
		}
	})

	b.Run("baseline--issue-credential", func(b usefulbench.Benchmarker) {
		for b.Loop() {
			params.GSK.IssueCredential(params.UPK)
		}
	})

	b.Run("baseline--sign", func(b usefulbench.Benchmarker) {
		for b.Loop() {
			_, err = zkban.Sign(params.M, params.CNT, params.Signer(), params.GPK, signCircuit.Prover())
			test.PanicIfErr(err)
		}
	})

	signature, _ := zkban.Sign(params.M, params.CNT, params.Signer(), params.GPK, signCircuit.Prover())

	b.Run("baseline--verify", func(b usefulbench.Benchmarker) {
		for b.Loop() {
			signature.Verify(params.M, params.CNT, params.Period, params.GPK, signCircuit.VerifyKey)
		}
	})

	if NoParallel {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
}

func benchmarkBasicUpdate(rl witness.RevocationList, name string, params *test.TestParams, b usefulbench.Benchmarker) {
	cacheName := "baseline-" + name
	prover, verifier := test.PrepareUpdateKeyCached(rl.Sizes(), cacheName)

	if NoParallel {
		runtime.GOMAXPROCS(1)
	}

	var err error
	b.Run("baseline--update-load-prove-"+name, func(b usefulbench.Benchmarker) {
		for b.Loop() {
			_, err = load.LoadUserKey(cacheName, "update")
			test.PanicIfErr(err)
		}
	})

	b.Run("baseline--update-load-verify-"+name, func(b usefulbench.Benchmarker) {
		for b.Loop() {
			_, err = load.LoadGroupManagerUpdateKey(cacheName)
			test.PanicIfErr(err)
		}
	})

	b.Run("baseline--update-req-"+name, func(b usefulbench.Benchmarker) {
		for b.Loop() {
			_, err := zkban.RequestUpdate(params.NextPeriod, params.Signer(), rl, params.GPK, prover)
			test.PanicIfErr(err)
		}
	})

	updateRequest, _ := zkban.RequestUpdate(params.NextPeriod, params.Signer(), rl, params.GPK, prover)

	vk, err := precomputes.NewUpdateVerificationKeyBLS12381(*verifier.VerifyKey)
	test.PanicIfErr(err)

	b.Run("baseline--update-precomputes-"+name, func(b usefulbench.Benchmarker) {
		for b.Loop() {
			_, err = vk.PrecomputeVerify(rl, params.GPK)
			test.PanicIfErr(err)
		}
	})

	prepared, _ := vk.PrecomputeVerify(rl, params.GPK)
	b.Run("baseline--update-verify-"+name, func(b usefulbench.Benchmarker) {
		for b.Loop() {
			err = vk.VerifyPrepared(*prepared, updateRequest, params.NextPeriod, params.Period)
			test.PanicIfErr(err)
		}
	})

	if NoParallel {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
}
