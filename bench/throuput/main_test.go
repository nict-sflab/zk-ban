package throuput

import (
	"runtime"
	"testing"

	zkban "github.com/akakou/zk-ban"
	"github.com/akakou/zk-ban/precomputes"
	"github.com/akakou/zk-ban/test"
	"github.com/akakou/zk-ban/witness"
)

func BenchmarkHandleParallel(b *testing.B) {
	b.Logf("GOMAXPROCS=%d", runtime.GOMAXPROCS(0))
	_, signCircuit, _ := test.PrepareCircuit(witness.EmptyRevocationList([]int{}))
	params := test.PrepareParams()

	signature, _ := zkban.Sign(params.M, params.CNT, params.Signer(), params.GPK, signCircuit.Prover())
	var err error

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err = signature.Verify(params.M, params.CNT, params.Period, params.GPK, signCircuit.VerifyKey)
			test.PanicIfErr(err)
		}
	})
	b.StopTimer()

	b.ReportMetric(float64(b.N)/b.Elapsed().Seconds(), "req/s")
}

func BenchmarkHandleParallel2(b *testing.B) {
	b.Logf("GOMAXPROCS=%d", runtime.GOMAXPROCS(0))
	params := test.PrepareParams()

	name := "throuput"
	cacheName := "baseline-" + name
	rl := test.EmptyUniformRevocationList(30, 30000)
	prover, verifier := test.PrepareUpdateKeyCached(rl.Sizes(), cacheName)

	var err error

	updateRequest, _ := zkban.RequestUpdate(params.NextPeriod, params.Signer(), rl, params.GPK, prover)
	vk, err := precomputes.NewUpdateVerificationKeyBLS12381(*verifier.VerifyKey)
	test.PanicIfErr(err)

	prepared, _ := vk.PrecomputeVerify(params.NextPeriod, params.Period, rl, params.GPK)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err = vk.VerifyPrepared(*prepared, updateRequest, params.NextPeriod, params.Period)
			test.PanicIfErr(err)

			_, err = params.GSK.IssueCredential(params.UPK)
			test.PanicIfErr(err)
		}
	})
	b.StopTimer()
	b.ReportMetric(float64(b.N)/b.Elapsed().Seconds(), "req/s")
}
