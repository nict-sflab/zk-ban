package circuit

import (
	"math/big"
	"testing"

	zkbanw "github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/test"
)

func TestSync(t *testing.T) {
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

	commitA, err := zkbanw.ComputeSignCommit(m, bsn, nonce, &usk)
	assert.NoError(err)

	commitB, err := zkbanw.ComputeSignCommit(m, bsn, nonce, &usk)
	assert.NoError(err)

	assign := &SyncCircuit{
		UserSecretKey: usk.Number,
		UserPublicKey: upk.Buffer,
		Period:        period,
		Commit2:       [RevocationListSize]frontend.Variable{commitA.Commit2, commitB.Commit2},
		Commit3:       [RevocationListSize]frontend.Variable{commitA.Commit3, commitB.Commit3},
	}

	assign.GroupPublicKey.Assign(twistededwards.BN254, gpk.Bytes())
	assign.Certificate.Assign(twistededwards.BN254, cert.Signature)

	assert.ProverSucceeded(&SyncCircuit{}, assign, test.WithCurves(ecc.BN254))
}
