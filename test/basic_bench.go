package test

import (
	"fmt"
	"testing"
	"time"

	gnarkprecomputes "github.com/akakou/gnark-precomputes"
	zkban "github.com/akakou/zk-ban"
	"github.com/akakou/zk-ban/circuit"
	zkbanc "github.com/akakou/zk-ban/circuit"
	"github.com/akakou/zk-ban/snark"
	bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/frontend"
)

func BenchmarkAll(t *testing.B) {
	params := prepareParams()

	rl1 := EmptyConstantRevocationAddList(270, 130)

	joinCircuit, signCircuit, updateCircuit1 := prepareCircuit(rl1, false)

	var proof groth16.Proof
	var err error

	var i time.Duration
	var n int

	var joinAssign *circuit.JoinRequestCircuit

	t.Run("join req", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			proof, joinAssign, err = zkban.JoinRequest(params.period, joinCircuit.Prover())
			panicIfErr(err)
		}
		i = b.Elapsed()
		n = b.N
	})

	fmt.Printf("zk-ban join bench: %v\n", i/time.Duration(n))

	t.Run("verify join req", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			wit, err := frontend.NewWitness(joinAssign, snark.EcCurve.ScalarField())
			panicIfErr(err)

			pubWit, err := wit.Public()
			panicIfErr(err)

			err = groth16.Verify(proof, joinCircuit.VerifyKey, pubWit)
			panicIfErr(err)
		}
	})

	var signAssign *circuit.SignCircuit
	t.Run("sign", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			proof, signAssign, err = zkban.Sign(params.m, params.cnt, params.signer(), params.gpk, signCircuit.Prover())
			panicIfErr(err)
		}

		i = b.Elapsed()
		n = b.N
	})

	fmt.Printf("zk-ban sign bench: %v\n", i/time.Duration(n))

	t.Run("verify", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			wit, err := frontend.NewWitness(signAssign, snark.EcCurve.ScalarField())
			panicIfErr(err)

			pubWit, err := wit.Public()
			panicIfErr(err)

			err = groth16.Verify(proof, signCircuit.VerifyKey, pubWit)
			panicIfErr(err)
		}
	})

	var updateAssign *circuit.UpdateCircuit
	var pubWit witness.Witness
	t.Run("update-req (constant)", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, proof, updateAssign, err = zkban.UpdateRequest(params.nextPeriod, params.signer(), rl1, params.gpk, updateCircuit1.Prover())
			panicIfErr(err)

			wit, err := frontend.NewWitness(updateAssign, snark.EcCurve.ScalarField())
			panicIfErr(err)

			pubWit, err = wit.Public()
			panicIfErr(err)
		}

		i = b.Elapsed()
		n = b.N
	})

	fmt.Printf("zk-ban update bench: %v\n", i/time.Duration(n))

	vk, err := gnarkprecomputes.FromBLS12381GnarkKey(updateCircuit1.VerifyKey, updateCircuit1.Circuit.(*zkbanc.UpdateCircuit))
	panicIfErr(err)

	var prepare *bls12381.G1Jac
	t.Run("update-verify-precomputes", func(b *testing.B) {
		for range b.N {
			prepare, err = vk.PreparePublicInputs(pubWit)
			panicIfErr(err)

		}
	})

	t.Run("update-verify (constant)", func(b *testing.B) {
		b.ResetTimer()

		for range b.N {
			params.gsk.IssueCredential(params.upk)
			err = vk.VerifyPrepared(proof, pubWit, prepare)
			panicIfErr(err)
		}
	})

}
