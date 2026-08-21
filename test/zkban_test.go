package test

import (
	"testing"

	zkban "github.com/akakou/zk-ban"
	"github.com/akakou/zk-ban/circuit"
	"github.com/akakou/zk-ban/precomputes"
	"github.com/akakou/zk-ban/primitives"
	zkbanw "github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/test"
)

func TestAll(t *testing.T) {
	assert := test.NewAssert(t)
	params := PrepareParams()

	proveFailedMessage := " is not satisfied:"
	verifyFailedMessage := "pairing doesn't match"

	rawRL := EmptyUniformRevocationList(2, 4)
	rl := AuthenticateRevocationList(params.GSK, rawRL)

	joinCircuit, signCircuit, updateCircuit := PrepareCircuit(rl)

	var signature *zkban.Signature
	var err error
	t.Run("join req", func(t *testing.T) {
		req, _, err := zkban.RequestJoin(params.Period, joinCircuit.Prover())
		assert.NoError(err)

		req.Verify(params.Period, joinCircuit.VerifyKey)
	})

	t.Run("sign", func(t *testing.T) {
		signature, err = zkban.Sign(params.M, params.CNT, params.Signer(), params.GPK, signCircuit.Prover())
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

		err = signature.Verify(primitives.NewBigInt(100000000), params.Period, params.GPK, signCircuit.VerifyKey)
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

	t.Run("update-revoked-fail", func(t *testing.T) {
		revokedRL := EmptyUniformRevocationList(2, 4)
		revokedRL[0].Period = primitives.NewBigInt(params.Period)
		revokedRL[0].Nyms[0] = signature.Commit.Nym
		revokedRL = AuthenticateRevocationList(params.GSK, revokedRL)

		_, err := zkban.RequestUpdate(params.NextPeriod, params.Signer(), revokedRL, params.GPK, updateCircuit.Prover())
		assert.ErrorContains(err, zkbanw.ErrRevokedNym.Error())
	})

	t.Run("update-invalid-interval-signature-fail", func(t *testing.T) {
		tamperedRL := tamperIntervalSignatures(rl)
		_, err := zkban.RequestUpdate(params.NextPeriod, params.Signer(), tamperedRL, params.GPK, updateCircuit.Prover())
		assert.ErrorContains(err, proveFailedMessage)
	})

	t.Run("update-signed-period-mismatch-fail", func(t *testing.T) {
		wrongPeriodRL := append(zkbanw.RevocationList(nil), rl...)
		wrongPeriodRL[0].Period = primitives.NewBigInt(rl[0].Period.Int64() + 1)

		_, err := zkban.RequestUpdate(params.NextPeriod, params.Signer(), wrongPeriodRL, params.GPK, updateCircuit.Prover())
		assert.ErrorContains(err, proveFailedMessage)
	})

	t.Run("update-public-period-vector-fail", func(t *testing.T) {
		req, err := zkban.RequestUpdate(params.NextPeriod, params.Signer(), rl, params.GPK, updateCircuit.Prover())
		assert.NoError(err)

		wrongPeriodRL := append(zkbanw.RevocationList(nil), rl...)
		wrongPeriodRL[0].Period = primitives.NewBigInt(rl[0].Period.Int64() + 1)

		err = req.Verify(params.NextPeriod, params.Period, wrongPeriodRL, params.GPK, updateCircuit.VerifyKey)
		assert.ErrorContains(err, verifyFailedMessage)

		vk, err := precomputes.NewUpdateVerificationKeyBLS12381(updateCircuit.VerifyKey)
		assert.NoError(err)
		prepared, err := vk.PrecomputeVerify(params.NextPeriod, params.Period, wrongPeriodRL, params.GPK)
		assert.NoError(err)

		err = vk.VerifyPrepared(*prepared, req, params.NextPeriod, params.Period)
		assert.ErrorContains(err, verifyFailedMessage)
	})
}

func tamperIntervalSignatures(rl zkbanw.RevocationList) zkbanw.RevocationList {
	result := make(zkbanw.RevocationList, len(rl))
	for periodIndex, revokedPerPeriod := range rl {
		result[periodIndex] = revokedPerPeriod
		result[periodIndex].SignedIntervals = make([]zkbanw.SignedInterval, len(revokedPerPeriod.SignedIntervals))
		for intervalIndex, interval := range revokedPerPeriod.SignedIntervals {
			result[periodIndex].SignedIntervals[intervalIndex] = interval
			result[periodIndex].SignedIntervals[intervalIndex].Signature = append([]byte(nil), interval.Signature...)
			if len(result[periodIndex].SignedIntervals[intervalIndex].Signature) > 0 {
				result[periodIndex].SignedIntervals[intervalIndex].Signature[len(result[periodIndex].SignedIntervals[intervalIndex].Signature)-1] ^= 1
			}
		}
	}
	return result
}
