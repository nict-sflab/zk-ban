package circuit

import (
	"math/big"
	"testing"

	"github.com/akakou/zk-ban/witness"
	"github.com/akakou/zk-ban/snark"
	"github.com/consensys/gnark/test"
)

func TestSignCircuit(t *testing.T) {
	assert := test.NewAssert(t)

	usk := witness.UserSecretKey{Number: big.NewInt(1)}
	m := big.NewInt(3)
	bsn := big.NewInt(4)
	period := big.NewInt(2024)

	gsk, gpk, err := witness.RandomGroupKeyPair()
	assert.NoError(err)

	upk, err := usk.PublicKey(period)
	assert.NoError(err)

	cert, err := gsk.IssueCredential(upk)
	assert.NoError(err)

	signer := witness.Signer{
		UserSecretKey: &usk,
		UserPublicKey: upk,
		Credential:    cert,
		Period:        period,
	}

	commit, err := signer.CommitSign(m, bsn)
	assert.NoError(err)

	authCircuit := SignCircuit{}

	witness := NewSignWitness(m, bsn, commit, &signer, gpk)

	assert.ProverSucceeded(&authCircuit, witness, test.WithCurves(snark.EcCurve))
}
