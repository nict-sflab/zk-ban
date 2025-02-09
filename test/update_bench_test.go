package zkbantest

import (
	"fmt"
	"testing"

	zkban "github.com/akakou/zk-ban"
	"github.com/akakou/zk-ban/commit"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
)

func BenchmarkUpdate(t *testing.B) {
	for i := 1; i <= 10; i++ {
		benchmarkUpdate(true, 270, 30*i, t)
	}

	for i := 1; i <= 10; i++ {
		benchmarkUpdate(false, 270, 270*30*i, t)
	}

	for i := 1; i <= 10; i++ {
		benchmarkUpdate(true, 1, 270*30*i, t)
	}

	for i := 1; i <= 10; i++ {
		benchmarkUpdate(true, 60*i, 35000/(60*i), t)
	}

	for i := 1; i <= 10; i++ {
		benchmarkUpdate(false, 60*i, 35000, t)
	}
}

func benchmarkUpdate(useConstant bool, a, b int, t *testing.B) {
	params := prepareParams()

	var proof groth16.Proof
	var pubWit witness.Witness
	var err error

	var rl commit.RevocationList
	if useConstant {
		rl = commit.EmptyConstantRevocationAddList(a, b)
	} else {
		rl = commit.EmptyLinerRevocationAddList(a, b)
	}

	_, _, updateCircuit := prepareCircuit(rl, true)

	tag1 := fmt.Sprintf("update-request (%v, %d, %d)", useConstant, a, b)
	tag2 := fmt.Sprintf("update-verify (%v, %d, %d)", useConstant, a, b)

	t.Run(tag1, func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, proof, pubWit, err = zkban.Update(params.nextPeriod, params.signer(), rl, params.gpk, updateCircuit.Prover())
			panicIfErr(err)
		}
	})

	t.Run(tag2, func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			params.gsk.IssueCredential(params.upk)

			err = groth16.Verify(proof, updateCircuit.VerifyKey, pubWit)
			panicIfErr(err)
		}
	})
}
