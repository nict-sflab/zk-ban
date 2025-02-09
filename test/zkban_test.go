package zkbantest

import (
	"encoding/json"
	"testing"

	zkban "github.com/akakou/zk-ban"
	"github.com/akakou/zk-ban/circuit"
	"github.com/akakou/zk-ban/commit"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/test"
)

func TestAll(t *testing.T) {
	assert := test.NewAssert(t)
	params := prepareParams()

	rl := commit.EmptyConstantRevocationAddList(270, 130)
	joinCircuit, signCircuit, updateCircuit := prepareCircuit(rl, false)

	t.Run("join req", func(t *testing.T) {
		var proof groth16.Proof
		proof, pubWit, _, _, err := zkban.JoinRequest(params.period, joinCircuit.Prover())
		assert.NoError(err)

		err = groth16.Verify(proof, joinCircuit.VerifyKey, pubWit)
		assert.NoError(err)
	})

	t.Run("sign", func(t *testing.T) {
		proof, pubWit, err := zkban.Sign(params.m, params.bsn, params.signer(), params.gpk, signCircuit.Prover())
		assert.NoError(err)

		err = groth16.Verify(proof, signCircuit.VerifyKey, pubWit)
		assert.NoError(err)
	})

	t.Run("sign-fail", func(t *testing.T) {
		proof, pubWit, err := zkban.Sign(params.m, params.bsn, params.signer(), params.gpk, signCircuit.Prover())
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

		err = groth16.Verify(proof, updateCircuit.VerifyKey, pubWit)
		assert.Error(err)
	})

	t.Run("update", func(t *testing.T) {
		_, proof, pubWit, err := zkban.Update(params.nextPeriod, params.signer(), rl, params.gpk, updateCircuit.Prover())
		assert.NoError(err)

		err = groth16.Verify(proof, updateCircuit.VerifyKey, pubWit)
		assert.NoError(err)
	})
}
