package test

import (
	"testing"

	zkban "github.com/akakou/zk-ban"
	"github.com/akakou/zk-ban/circuit"
	"github.com/akakou/zk-ban/precomputes"
	"github.com/akakou/zk-ban/primitives"
	"github.com/consensys/gnark/test"
)

func TestAll(t *testing.T) {
	assert := test.NewAssert(t)
	params := PrepareParams()

	proveFailedMessage := " is not satisfied:"
	verifyFailedMessage := "pairing doesn't match"

	rl := EmptyUniformRevocationList(60, 120)

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

		err = signature.Verify(params.M, params.Period, params.GPK, signCircuit.VerifyKey)
		assert.NoError(err)
	})

	t.Run("sign-precomputes", func(t *testing.T) {
		signature, err := zkban.Sign(params.M, params.CNT, params.Signer(), params.GPK, signCircuit.Prover())
		assert.NoError(err)

		vk, err := precomputes.NewAuthVerificationKeyBLS12381(signCircuit.VerifyKey)
		assert.NoError(err)
		prepared, err := vk.PrecomputeVerify(params.Period, params.GPK)
		assert.NoError(err)

		err = vk.VerifyPrepared(*prepared, params.M, signature)
		assert.NoError(err)
	})

	t.Run("sign-fail", func(t *testing.T) {
		signature, err = zkban.Sign(params.M, params.CNT, params.Signer(), params.GPK, signCircuit.Prover())
		assert.NoError(err)

		err = signature.Verify(primitives.NewBigInt(100000000), params.CNT, params.Period, params.GPK, signCircuit.VerifyKey)
		assert.ErrorContains(err, verifyFailedMessage)
	})

	t.Run("sign-counter-max-fail", func(t *testing.T) {
		_, err = zkban.Sign(params.M, int64(circuit.MaxSession), params.Signer(), params.GPK, signCircuit.Prover())
		assert.ErrorContains(err, proveFailedMessage)
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

		prepared, err := vk.PrecomputeVerify(params.NextPeriod, params.Period, rl, params.GPK)
		assert.NoError(err)

		err = vk.VerifyPrepared(*prepared, req, params.NextPeriod, params.Period)
		assert.NoError(err)
	})

	t.Run("update-fail1", func(t *testing.T) {
		rl := EmptyUniformRevocationList(60, 120)
		rl[0].Period = primitives.NewBigInt(params.Period)
		rl[0].Nyms[0] = signature.Commit.Nym
		_, err := zkban.RequestUpdate(params.NextPeriod, params.Signer(), rl, params.GPK, updateCircuit.Prover())
		assert.ErrorContains(err, proveFailedMessage)
	})

	t.Run("update-fail2", func(t *testing.T) {
		test := func(periodIndex, nymIndex int) {
			rl := EmptyUniformRevocationList(60, 120)
			req, err := zkban.RequestUpdate(params.NextPeriod, params.Signer(), rl, params.GPK, updateCircuit.Prover())
			assert.NoError(err, proveFailedMessage)

			rl[periodIndex].Period = primitives.NewBigInt(params.Period)
			rl[periodIndex].Nyms[nymIndex] = signature.Commit.Nym

			vk, err := precomputes.NewUpdateVerificationKeyBLS12381(updateCircuit.VerifyKey)
			assert.NoError(err, proveFailedMessage)

			cache, err := vk.PrecomputeVerify(params.NextPeriod, params.Period, rl, params.GPK)
			assert.NoError(err, proveFailedMessage)

			err = vk.VerifyPrepared(*cache, req, params.NextPeriod, params.Period)
			assert.ErrorContains(err, verifyFailedMessage)
		}

		test(0, 0)
		test(59, 1)
	})
}
