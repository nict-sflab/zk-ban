package circuit

import (
	"math/big"
	"testing"

	"github.com/akakou/zk-ban/snark"
	zkbanw "github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/test"
)

func TestSignCircuit(t *testing.T) {
	assert := test.NewAssert(t)

	usk := zkbanw.UserSecretKey{Number: big.NewInt(1)}
	nonce := big.NewInt(2)
	m := big.NewInt(3)
	bsn := big.NewInt(4)
	period := big.NewInt(2024)

	gsk, gpk, err := zkbanw.RandomGroupKeyPair()
	assert.NoError(err)

	upk, err := usk.PublicKey(period)
	assert.NoError(err)

	cert, err := gsk.IssueCertificate(upk)
	assert.NoError(err)

	commit, err := zkbanw.ComputeSignCommit(m, bsn, nonce, &usk)
	assert.NoError(err)

	authCircuit := SignCircuit{}

	assign := &SignCircuit{
		UserPublicKey: upk.Buffer,
		Nonce:         nonce,
		UserSecretKey: usk.Number,
		Basename:      bsn,
		Message:       m,
		Commit1:       commit.Commit1,
		Commit2:       commit.Commit2,
		Commit3:       commit.Commit3,
		Period:        period,
	}

	assign.GroupPublicKey.Assign(snark.TwistededwardsCurve, gpk.Bytes())
	assign.Certificate.Assign(snark.TwistededwardsCurve, cert.Signature)

	assert.ProverSucceeded(&authCircuit, assign, test.WithCurves(snark.EcCurve))
}
