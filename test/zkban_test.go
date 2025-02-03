package zkban_test

import (
	"encoding/json"
	"testing"

	zkban "github.com/akakou/zk-ban"
	"github.com/akakou/zk-ban/circuit"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
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

	t.Run("sign-fail", func(t *testing.T) {
		proof, pubWit, err := zkban.Sign(params.m, params.bsn, params.signer(), params.gpk, params.signSnark.Prover())
		assert.NoError(err)

		schema, err := frontend.NewSchema(&circuit.SignCircuit{})
		assert.NoError(err)

		data, err := pubWit.ToJSON(schema)
		assert.NoError(err)

		circuit := circuit.SignCircuit{}
		err = json.Unmarshal(data, &circuit)
		assert.NoError(err)

		circuit.Message = 20000

		data, err = json.Marshal(&circuit)
		assert.NoError(err)

		err = pubWit.FromJSON(schema, data)
		assert.NoError(err)

		err = groth16.Verify(proof, params.signSnark.VerifyKey, pubWit)
		assert.Error(err)
	})

	t.Run("update", func(t *testing.T) {
		_, proof, pubWit, err := zkban.Update(params.nextPeriod, params.signer(), params.revocationList, params.gpk, params.updateSnark.Prover())
		assert.NoError(err)

		err = groth16.Verify(proof, params.updateSnark.VerifyKey, pubWit)
		assert.NoError(err)
	})
}
