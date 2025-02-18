package zkbantest

import (
	"encoding/json"
	"math/big"
	"testing"

	zkban "github.com/akakou/zk-ban"
	"github.com/akakou/zk-ban/circuit"
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
		proof, pubWit, _, _, err := zkban.JoinRequest(params.period, joinCircuit.Prover())
		assert.NoError(err)

		err = groth16.Verify(proof, joinCircuit.VerifyKey, pubWit)
		assert.NoError(err)
	})

	t.Run("sign", func(t *testing.T) {
		proof, pubWit, err := zkban.Sign(params.m, params.cnt, params.signer(), params.gpk, signCircuit.Prover())
		assert.NoError(err)

		err = groth16.Verify(proof, signCircuit.VerifyKey, pubWit)
		assert.NoError(err)
	})

	sig := circuit.SignCircuit{}
	t.Run("sign-fail", func(t *testing.T) {
		proof, pubWit, err := zkban.Sign(params.m, params.cnt, params.signer(), params.gpk, signCircuit.Prover())
		assert.NoError(err)

		schema, err := frontend.NewSchema(&circuit.SignCircuit{})
		assert.NoError(err)

		data, err := pubWit.ToJSON(schema)
		assert.NoError(err)

		err = json.Unmarshal(data, &sig)
		assert.NoError(err)

		tmp := sig.SessionTag
		sig.SessionTag = 20000

		data, err = json.Marshal(&sig)
		assert.NoError(err)

		err = pubWit.FromJSON(schema, data)
		assert.NoError(err)

		err = groth16.Verify(proof, signCircuit.VerifyKey, pubWit)
		assert.ErrorContains(err, verifyFailedMessage)

		// used by update-fail
		sig.SessionTag = tmp
	})

	t.Run("update", func(t *testing.T) {
		_, proof, pubWit, err := zkban.Update(params.nextPeriod, params.gpk, params.signer(), rl, updateCircuit.Prover())
		assert.NoError(err)

		err = groth16.Verify(proof, updateCircuit.VerifyKey, pubWit)
		assert.NoError(err)
	})

	t.Run("update-fail", func(t *testing.T) {
		rl[0].SessionTag = big.NewInt(int64(sig.SessionTag.(float64)))
		rl[0].Nyms[0] = big.NewInt(0)
		rl[0].Nyms[0].SetString(sig.Nym.(string), 10)

		_, _, _, err := zkban.Update(params.nextPeriod, params.gpk, params.signer(), rl, updateCircuit.Prover())
		assert.ErrorContains(err, proveFailedMessage)

		rl[0].Nyms[0] = big.NewInt(0)

		_, proof, pubWit, err := zkban.Update(params.nextPeriod, params.gpk, params.signer(), rl, updateCircuit.Prover())
		assert.NoError(err)

		witnessRL := circuit.NewRevocationListWitness(rldash)
		schema, err := frontend.NewSchema(&circuit.UpdateCircuit{
			RevocationList: witnessRL,
		})
		assert.NoError(err)

		data, err := pubWit.ToJSON(schema)
		assert.NoError(err)

		update := circuit.UpdateCircuit{}
		err = json.Unmarshal(data, &update)
		assert.NoError(err)

		update.RevocationList[0].Nyms[0] = sig.Nym

		data, err = json.Marshal(&update)
		assert.NoError(err)

		err = pubWit.FromJSON(schema, data)
		assert.NoError(err)

		err = groth16.Verify(proof, updateCircuit.VerifyKey, pubWit)
		assert.ErrorContains(err, verifyFailedMessage)
	})
}
