package zkbantest

import (
	"fmt"
	"testing"

	zkban "github.com/akakou/zk-ban"
	zkbanw "github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
)

func BenchmarkUpdate(t *testing.B) {
	max := 10
	base := 600
	session := 270
	revocation := 35100

	for i := 1; i <= max; i++ {
		benchmarkUpdate(true, session, base*i, t)
	}

	for i := 1; i <= max; i++ {
		benchmarkUpdate(false, session, session*base*i, t)
	}

	for i := 1; i <= max; i++ {
		benchmarkUpdate(true, 1, session*base*i, t)
	}

	for i := 1; i <= max; i++ {
		benchmarkUpdate(true, base*i, revocation/(base*i), t)
	}

	for i := 1; i <= max; i++ {
		benchmarkUpdate(false, base*i, revocation, t)
	}
}

func benchmarkUpdate(useConstant bool, a, b int, t *testing.B) {
	params := prepareParams()

	var proof groth16.Proof
	var pubWit witness.Witness
	var err error

	var rl zkbanw.RevocationList
	var tag = ""
	if useConstant {
		rl = zkbanw.EmptyConstantRevocationAddList(a, b)
		tag = fmt.Sprintf("%d,%d,%d,%v", a*b, a, b, "constant")

	} else {
		rl = zkbanw.EmptyLinerRevocationAddList(a, b)
		tag = fmt.Sprintf("%d,%d,%v,%v", b, a, "-", "linear")
	}

	_, _, updateCircuit := prepareCircuit(rl, true)

	t.Run("Prove ,"+tag, func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, proof, pubWit, err = zkban.Update(params.nextPeriod, params.gpk, params.signer(), rl, updateCircuit.Prover())
			panicIfErr(err)
		}
	})

	t.Run("Verify ,"+tag, func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			params.gsk.IssueCredential(params.upk)

			err = groth16.Verify(proof, updateCircuit.VerifyKey, pubWit)
			panicIfErr(err)
		}
	})
}
