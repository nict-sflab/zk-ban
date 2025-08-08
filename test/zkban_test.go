package zkbantest

import (
	"testing"

	zkban "github.com/akakou/zk-ban"
	"github.com/akakou/zk-ban/precomputes"
	"github.com/akakou/zk-ban/primitives"
	"github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/test"
)

func TestAll(t *testing.T) {
	assert := test.NewAssert(t)
	params := prepareParams()

	proveFailedMessage := " is not satisfied:"
	verifyFailedMessage := "pairing doesn't match"

	rl := witness.EmptyConstantRevocationAddList(270, 130)

	joinCircuit, signCircuit, updateCircuit := prepareCircuit(rl, false)

	var signature *zkban.Signature
	var err error
	t.Run("join req", func(t *testing.T) {
		req, _, err := zkban.RequestJoin(params.period, joinCircuit.Prover())
		assert.NoError(err)

		req.Verify(params.period, joinCircuit.VerifyKey)
	})

	t.Run("sign", func(t *testing.T) {
		signature, err := zkban.Sign(params.m, params.cnt, params.signer(), params.gpk, signCircuit.Prover())
		assert.NoError(err)

		err = signature.Verify(params.m, params.cnt, params.period, params.gpk, signCircuit.VerifyKey)
		assert.NoError(err)
	})

	t.Run("sign-fail", func(t *testing.T) {
		signature, err = zkban.Sign(params.m, params.cnt, params.signer(), params.gpk, signCircuit.Prover())
		assert.NoError(err)

		err = signature.Verify(primitives.NewBigInt(100000000), params.cnt, params.period, params.gpk, signCircuit.VerifyKey)
		assert.ErrorContains(err, verifyFailedMessage)
	})

	t.Run("update", func(t *testing.T) {
		req, err := zkban.RequestUpdate(params.nextPeriod, params.signer(), rl, params.gpk, updateCircuit.Prover())
		assert.NoError(err)

		err = req.Verify(params.nextPeriod, params.period, rl, params.gpk, updateCircuit.VerifyKey)
		assert.NoError(err)
	})

	t.Run("update-precomputes", func(t *testing.T) {
		req, err := zkban.RequestUpdate(params.nextPeriod, params.signer(), rl, params.gpk, updateCircuit.Prover())
		assert.NoError(err)

		vk, err := precomputes.NewUpdateVerificationKeyBLS12381(updateCircuit.VerifyKey)
		assert.NoError(err)

		prepared, err := vk.PrecomputeVerify(req, params.nextPeriod, params.period, rl, params.gpk)
		assert.NoError(err)

		err = vk.VerifyPrepared(*prepared, req, params.nextPeriod, params.period, rl, params.gsk, params.gpk)
		assert.NoError(err)
	})

	t.Run("update-fail", func(t *testing.T) {
		rl[0].SessionTag = &witness.SessionTag(params.cnt, params.period).Int
		rl[0].Nyms[0] = &signature.Commit.Nym.Int

		_, err := zkban.RequestUpdate(params.nextPeriod, params.signer(), rl, params.gpk, updateCircuit.Prover())
		assert.ErrorContains(err, proveFailedMessage)

		// err = req.Verify(params.nextPeriod, params.period, rl, params.gpk, updateCircuit.VerifyKey)
		// assert.ErrorContains(err, verifyFailedMessage)

		// _, _, _, err := zkban.RequestUpdate(params.nextPeriod, params.signer(), rl, params.gpk, updateCircuit.Prover())
		// assert.ErrorContains(err, proveFailedMessage)

		// rl[0].Nyms[0] = big.NewInt(0)

		// _, proof, assign, err := zkban.RequestUpdate(params.nextPeriod, params.signer(), rl, params.gpk, updateCircuit.Prover())
		// assert.NoError(err)

		// witnessRL := circuit.NewRevocationListWitness(rldash)
		// assert.NoError(err)

		// assign.RevocationList = witnessRL

		// wit, err := frontend.NewWitness(assign, snark.EcCurve.ScalarField())
		// assert.NoError(err)

		// pubWit, err := wit.Public()
		// assert.NoError(err)

		// err = groth16.Verify(proof, updateCircuit.VerifyKey, pubWit)
		// assert.ErrorContains(err, verifyFailedMessage)
	})
}
