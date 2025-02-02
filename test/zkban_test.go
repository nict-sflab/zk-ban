package zkban_test

import (
	"testing"

	zkban "github.com/akakou/zk-ban"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/test"
)

func TestAll(t *testing.T) {
	assert := test.NewAssert(t)
	params := prepare()

	t.Run("join req", func(t *testing.T) {
		var proof groth16.Proof
		proof, pubWit, _, _, err := zkban.JoinRequest(params.period, params.joinSnark.Prover())
		assert.NoError(err)

		err = groth16.Verify(proof, params.joinSnark.VerifyKey, pubWit)
		assert.NoError(err)
	})

	t.Run("sign", func(t *testing.T) {
		proof, pubWit, err := zkban.Sign(params.m, params.bsn, params.signer(), params.gpk, params.signSnark.Prover())
		assert.NoError(err)

		err = groth16.Verify(proof, params.signSnark.VerifyKey, pubWit)
		assert.NoError(err)
	})

	t.Run("update", func(t *testing.T) {
		_, proof, pubWit, err := zkban.Update(params.revocationList, params.nextPeriod, params.signer(), params.sessionName, params.gpk, params.updateSnark.Prover())
		assert.NoError(err)

		err = groth16.Verify(proof, params.updateSnark.VerifyKey, pubWit)
		assert.NoError(err)
	})
}
