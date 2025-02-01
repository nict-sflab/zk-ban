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
		UserSecretKey: &usk,
		UserPublicKey: upk,
		Credential:    cert,
		Period:        period,
	}

	commit, err := signer.ComputeSignCommit(m, bsn)
	assert.NoError(err)

	authCircuit := SignCircuit{}

	assign := &SignCircuit{
		UserSecretKey: usk.Number,
		Basename:      bsn,
		Message:       m,
		Commit1:       commit.Commit1,
		Commit2:       commit.Commit2,
		Period:        period,
	}

	assign.GroupPublicKey.Assign(snark.TwistededwardsCurve, gpk.Bytes())
	assign.Credential.Assign(snark.TwistededwardsCurve, cert.Signature)

	assert.ProverSucceeded(&authCircuit, assign, test.WithCurves(snark.EcCurve))
}
