package circuit

import (
	"math/big"
	"testing"

	"github.com/akakou/zk-ban/commit"
	"github.com/akakou/zk-ban/snark"
	"github.com/consensys/gnark/test"
)

func TestSignCircuit(t *testing.T) {
	assert := test.NewAssert(t)

	usk := commit.UserSecretKey{Number: big.NewInt(1)}
	m := big.NewInt(3)
	bsn := big.NewInt(4)
	period := big.NewInt(2024)

	gsk, gpk, err := commit.RandomGroupKeyPair()
	assert.NoError(err)

	upk, err := usk.PublicKey(period)
	assert.NoError(err)

	cert, err := gsk.IssueCredential(upk)
	assert.NoError(err)

	signer := commit.Signer{
		UserSecretKey:  &usk,
		UserPublicKey:  upk,
		Credential:     cert,
		Period:         period,
		GroupPublicKey: gpk,
	}

	commit, err := signer.CommitSign(m, bsn)
	assert.NoError(err)

	authCircuit := SignCircuit{}

	witness := NewSignWitness(m, bsn, commit, &signer)

	assert.ProverSucceeded(&authCircuit, witness, test.WithCurves(snark.EcCurve))
}
