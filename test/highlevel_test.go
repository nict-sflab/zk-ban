package zkbantest

import (
	"testing"

	"github.com/akakou/zk-ban/highlevel"
	"github.com/akakou/zk-ban/snark"
	zkbanw "github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/test"
)

func TestHighLevelApi(t *testing.T) {
	params := prepareParams()
	assert := test.NewAssert(t)

	rl1 := zkbanw.EmptyConstantRevocationAddList(10, 10)
	joinCircuit, signCircuit, updateCircuit := prepareCircuit(rl1, false)

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

	updateCircuitBytes, err := snark.EncodeCircuit(updateCircuit.ConstraintSystem)
	assert.NoError(err)

	updateProveKeyBytes, err := snark.EncodeProverKey(updateCircuit.ProveKey)
	assert.NoError(err)

	updateVerifyKeyBytes, err := snark.EncodeVerifierKey(updateCircuit.VerifyKey)
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

	cred, err := highlevel.IssueCredential(params.period.Int64(), req.UserPublicKey, params.gsk.Bytes())
	assert.NoError(err)

	signer := highlevel.HighLevelSigner{
		Credential:     cred,
		GroupPublicKey: params.gpk.Bytes(),
		Period:         params.signer().Period.Int64(),
		Secret:         req.UserSecretKey,
		UserPublicKey:  req.UserPublicKey,
	}

	signature, err := highlevel.Sign(
		message,
		params.cnt.Int64(),
		signer,
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

	updateSigner, proof, err := highlevel.UpdateRequest(
		params.nextPeriod.Int64(),
		&signer,
		rl1,
		params.gpk.Bytes(),
		&highlevel.HighLevelSnarkProver{
			ConstraintSystem: updateCircuitBytes,
			ProveKey:         updateProveKeyBytes,
		},
	)

	assert.NoError(err)

	err = highlevel.VerifyUpdateRequest(
		proof,
		updateSigner.UserPublicKey,
		params.nextPeriod.Int64(),
		&signer,
		rl1,
		updateVerifyKeyBytes,
	)

	assert.NoError(err)
}
