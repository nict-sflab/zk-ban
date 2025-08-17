package test

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
	params := PrepareParams()

	proveFailedMessage := " is not satisfied:"
	verifyFailedMessage := "pairing doesn't match"

	rl := EmptyUniformRevocationList(180, 300)

	joinCircuit, signCircuit, updateCircuit := PrepareCircuit(rl)

	var signature *zkban.Signature
	var err error
	t.Run("join req", func(t *testing.T) {
		req, _, err := zkban.RequestJoin(params.Period, joinCircuit.Prover())
		assert.NoError(err)

		req.Verify(params.Period, joinCircuit.VerifyKey)
	})

	t.Run("sign", func(t *testing.T) {
		signature, err := zkban.Sign(params.M, params.CNT, params.Signer(), params.GPK, signCircuit.Prover())
		assert.NoError(err)

		err = signature.Verify(params.M, params.CNT, params.Period, params.GPK, signCircuit.VerifyKey)
		assert.NoError(err)
	})

	t.Run("sign-fail", func(t *testing.T) {
		signature, err = zkban.Sign(params.M, params.CNT, params.Signer(), params.GPK, signCircuit.Prover())
		assert.NoError(err)

		err = signature.Verify(primitives.NewBigInt(100000000), params.CNT, params.Period, params.GPK, signCircuit.VerifyKey)
		assert.ErrorContains(err, verifyFailedMessage)
	})

	t.Run("update", func(t *testing.T) {
		req, err := zkban.RequestUpdate(params.NextPeriod, params.Signer(), rl, params.GPK, updateCircuit.Prover())
		assert.NoError(err)

		err = req.Verify(params.NextPeriod, params.Period, rl, params.GPK, updateCircuit.VerifyKey)
		assert.NoError(err)
	})

	t.Run("update-precomputes", func(t *testing.T) {
		req, err := zkban.RequestUpdate(params.NextPeriod, params.Signer(), rl, params.GPK, updateCircuit.Prover())
		assert.NoError(err)

		vk, err := precomputes.NewUpdateVerificationKeyBLS12381(updateCircuit.VerifyKey)
		assert.NoError(err)

		prepared, err := vk.PrecomputeVerify(rl, params.GPK)
		assert.NoError(err)

		err = vk.VerifyPrepared(*prepared, req, params.NextPeriod, params.Period)
		assert.NoError(err)
	})

	t.Run("update-fail", func(t *testing.T) {
		rl[0].SessionTag = witness.SessionTag(params.CNT, params.Period)
		rl[0].Nyms[0] = signature.Commit.Nym

		_, err := zkban.RequestUpdate(params.NextPeriod, params.Signer(), rl, params.GPK, updateCircuit.Prover())
		assert.ErrorContains(err, proveFailedMessage)

		// err = req.Verify(params.nextPeriod, params.period, rl, params.gpk, updateCircuit.VerifyKey)
		// assert.ErrorContains(err, verifyFailedMessage)

		// _, _, _, err := zkban.RequestUpdate(params.nextPeriod, params.Signer(), rl, params.gpk, updateCircuit.Prover())
		// assert.ErrorContains(err, proveFailedMessage)

		// rl[0].Nyms[0] = big.NewInt(0)

		// _, proof, assign, err := zkban.RequestUpdate(params.nextPeriod, params.Signer(), rl, params.gpk, updateCircuit.Prover())
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
