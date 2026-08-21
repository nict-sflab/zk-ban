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

var BaseSubRevocation = 748
var BasePeriods = 13

var Name = ""

func BenchmarkBaseline(b usefulbench.Benchmarker) {
	params := test.PrepareParams()

	ub, ok := b.(*usefulbench.UsefulBenchmaker)

	if ok {
		ub.Result["env"]["basePeriodNum"] = int64(BasePeriods)
		ub.Result["env"]["baseNymNum"] = int64(BaseSubRevocation)
	}

	if !OnlyUpdate {
		benchmarkBaseline(&params, b)
	}

	rlU := test.EmptyUniformRevocationList(BasePeriods, BaseSubRevocation)
	benchmarkBasicUpdate(rlU, Name+"uniform", &params, b)

	if OnlyUniform {
		return
	}

	rlP := test.EmptyProportionalRevocationList(BasePeriods, BaseSubRevocation)
	benchmarkBasicUpdate(rlP, Name+"proportional", &params, b)

	rlG := test.EmptyGaussianRevocationList(BasePeriods, BaseSubRevocation)
	benchmarkBasicUpdate(rlG, Name+"gaussian", &params, b)
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
	vk, err := precomputes.NewAuthVerificationKeyBLS12381(signCircuit.VerifyKey)
	test.PanicIfErr(err)

	prepared, err := vk.PrecomputeVerify(params.Period, params.GPK)
	test.PanicIfErr(err)

	test.PanicIfErr(err)

	b.Run("baseline--verify", func(b usefulbench.Benchmarker) {
		for b.Loop() {
			err = vk.VerifyPrepared(*prepared, params.M, signature)
		}
	})

	if NoParallel {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
}

func benchmarkBasicUpdate(rl witness.RevocationList, name string, params *test.TestParams, b usefulbench.Benchmarker) {
	rl = test.AuthenticateRevocationList(params.GSK, rl)
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
			_, err = vk.PrecomputeVerify(params.NextPeriod, params.Period, rl, params.GPK)
			test.PanicIfErr(err)
		}
	})

	prepared, _ := vk.PrecomputeVerify(params.NextPeriod, params.Period, rl, params.GPK)
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
