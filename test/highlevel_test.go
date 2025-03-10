package zkbantest

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/akakou/zk-ban/highlevel"
	zkbanw "github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/test"
)

func TestHighLevelApi(t *testing.T) {
	params := prepareParams()
	assert := test.NewAssert(t)

	rl1 := zkbanw.EmptyConstantRevocationAddList(1, 1)
	_, signCircuit, _ := prepareCircuit(rl1, false)

	var circuitBuf bytes.Buffer
	_, err := signCircuit.ConstraintSystem.WriteTo(&circuitBuf)
	assert.NoError(err)
	circuitBytes := circuitBuf.Bytes()

	var proveKeyBuf bytes.Buffer
	err = signCircuit.ProveKey.WriteDump(&proveKeyBuf)
	assert.NoError(err)
	proveKeyBytes := proveKeyBuf.Bytes()

	var verifyKeyBuf bytes.Buffer
	_, err = signCircuit.VerifyKey.WriteTo(&verifyKeyBuf)
	assert.NoError(err)
	verifyKeyBytes := verifyKeyBuf.Bytes()

	message := params.m.Bytes()

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
		circuitBytes,
		proveKeyBytes,
	)

	assert.NoError(err)

	err = highlevel.Verify(
		signature,
		params.m.Bytes(),
		params.cnt.Int64(),
		params.period.Int64(),
		params.gpk.Bytes(),
		circuitBytes,
		verifyKeyBytes,
	)

	assert.NoError(err)
}
