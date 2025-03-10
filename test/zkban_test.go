package zkbantest

import (
	"math/big"
	"testing"

	zkban "github.com/akakou/zk-ban"
	"github.com/akakou/zk-ban/circuit"
	"github.com/akakou/zk-ban/snark"
	"github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/test"
)

func TestAll(t *testing.T) {
	assert := test.NewAssert(t)
	params := prepareParams()

	proveFailedMessage := " is not satisfied:"
	verifyFailedMessage := "pairing doesn't match"

	rl := witness.EmptyConstantRevocationAddList(270, 130)
	rldash := witness.EmptyConstantRevocationAddList(270, 130)

	joinCircuit, signCircuit, updateCircuit := prepareCircuit(rl, false)

	t.Run("join req", func(t *testing.T) {
		var proof groth16.Proof
		proof, assign, err := zkban.JoinRequest(params.period, joinCircuit.Prover())
		assert.NoError(err)

		wit, err := frontend.NewWitness(assign, snark.EcCurve.ScalarField())
		assert.NoError(err)

		pubWit, err := wit.Public()
		assert.NoError(err)

		err = groth16.Verify(proof, joinCircuit.VerifyKey, pubWit)
		assert.NoError(err)
	})

	var sign *circuit.SignCircuit
	t.Run("sign", func(t *testing.T) {
		proof, assign, err := zkban.Sign(params.m, params.cnt, params.signer(), params.gpk, signCircuit.Prover())
		assert.NoError(err)

		sign = assign

		wit, err := frontend.NewWitness(assign, snark.EcCurve.ScalarField())
		assert.NoError(err)

		pubWit, err := wit.Public()
		assert.NoError(err)

		err = groth16.Verify(proof, signCircuit.VerifyKey, pubWit)
		assert.NoError(err)
	})

	t.Run("sign-fail", func(t *testing.T) {
		proof, assign, err := zkban.Sign(params.m, params.cnt, params.signer(), params.gpk, signCircuit.Prover())
		assert.NoError(err)

		assign.Message = big.NewInt(100000)
		wit, err := frontend.NewWitness(assign, snark.EcCurve.ScalarField())
		assert.NoError(err)

		pubWit, err := wit.Public()
		assert.NoError(err)

		err = groth16.Verify(proof, signCircuit.VerifyKey, pubWit)
		assert.ErrorContains(err, verifyFailedMessage)
	})

	t.Run("update", func(t *testing.T) {
		_, proof, assign, err := zkban.UpdateRequest(params.nextPeriod, params.signer(), rl, params.gpk, updateCircuit.Prover())
		assert.NoError(err)

		wit, err := frontend.NewWitness(assign, snark.EcCurve.ScalarField())
		assert.NoError(err)

		pubWit, err := wit.Public()
		assert.NoError(err)

		err = groth16.Verify(proof, updateCircuit.VerifyKey, pubWit)
		assert.NoError(err)
	})

	t.Run("update-fail", func(t *testing.T) {
		rl[0].SessionTag = sign.SessionTag.(*big.Int)
		rl[0].Nyms[0] = sign.Nym.(*big.Int)

		_, _, _, err := zkban.UpdateRequest(params.nextPeriod, params.signer(), rl, params.gpk, updateCircuit.Prover())
		assert.ErrorContains(err, proveFailedMessage)

		rl[0].Nyms[0] = big.NewInt(0)

		_, proof, assign, err := zkban.UpdateRequest(params.nextPeriod, params.signer(), rl, params.gpk, updateCircuit.Prover())
		assert.NoError(err)

		witnessRL := circuit.NewRevocationListWitness(rldash)
		assert.NoError(err)

		assign.RevocationList = witnessRL

		wit, err := frontend.NewWitness(assign, snark.EcCurve.ScalarField())
		assert.NoError(err)

		pubWit, err := wit.Public()
		assert.NoError(err)

		err = groth16.Verify(proof, updateCircuit.VerifyKey, pubWit)
		assert.ErrorContains(err, verifyFailedMessage)
	})
}
