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

	result := highlevel.Sign(
		message,
		params.cnt.Int64(),
		signerBytes,
		circuitBytes,
		proveKeyBytes,
	)

	var resultStrct highlevel.Result
	err = resultStrct.FromString(result)
	assert.NoError(err)

	err = json.Unmarshal([]byte(result), &resultStrct)
	assert.NoError(err)
	assert.Empty(resultStrct.Err)

	result = highlevel.Verify(
		resultStrct.Out,
		params.m.Bytes(),
		params.cnt.Int64(),
		params.period.Int64(),
		params.gpk.Bytes(),
		circuitBytes,
		verifyKeyBytes,
	)

	err = resultStrct.FromString(result)
	assert.NoError(err)
	assert.Equal(resultStrct.Err, "")
}
