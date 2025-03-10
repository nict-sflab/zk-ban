package zkbantest

import (
	"encoding/json"
	"testing"

	"github.com/akakou/zk-ban/highlevel"
	"github.com/akakou/zk-ban/snark"
	zkbanw "github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/test"
)

func TestHighLevelApi(t *testing.T) {
	params := prepareParams()
	assert := test.NewAssert(t)

	rl1 := zkbanw.EmptyConstantRevocationAddList(1, 1)
	joinCircuit, signCircuit, _ := prepareCircuit(rl1, false)

	joinCircuitBytes, err := snark.EncodeCircuit(joinCircuit.ConstraintSystem)
	assert.NoError(err)

	joinProveKeyBytes, err := snark.EncodeProverKey(joinCircuit.ProveKey)
	assert.NoError(err)

	joinVerifyKeyBytes, err := snark.EncodeVerifierKey(joinCircuit.VerifyKey)
	assert.NoError(err)

	signCircuitBytes, err := snark.EncodeCircuit(signCircuit.ConstraintSystem)
	assert.NoError(err)

	signProveKeyBytes, err := snark.EncodeProverKey(signCircuit.ProveKey)
	assert.NoError(err)

	signVerifyKeyBytes, err := snark.EncodeVerifierKey(signCircuit.VerifyKey)
	assert.NoError(err)

	message := params.m.Bytes()

	proof, req, err := highlevel.JoinRequest(
		params.period.Int64(),
		joinCircuitBytes,
		joinProveKeyBytes,
	)

	assert.NoError(err)

	err = highlevel.VerifyJoinReq(
		proof,
		req.UserPublicKey,
		req.Period,
		joinVerifyKeyBytes,
	)

	assert.NoError(err)

	signer := highlevel.Signer{
		Credential:     params.cert.Signature,
		GroupPublicKey: params.gpk.Bytes(),
		Period:         params.signer().Period.Int64(),
		Secret:         params.usk.Number.Bytes(),
		UserPublicKey:  params.upk.Number.Bytes(),
	}

	signerBytes, err := json.Marshal(signer)
	assert.NoError(err)

	signature, err := highlevel.Sign(
		message,
		params.cnt.Int64(),
		signerBytes,
		signCircuitBytes,
		signProveKeyBytes,
	)

	assert.NoError(err)

	err = highlevel.Verify(
		signature,
		params.m.Bytes(),
		params.cnt.Int64(),
		params.period.Int64(),
		params.gpk.Bytes(),
		signVerifyKeyBytes,
	)

	assert.NoError(err)
}
