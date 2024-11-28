package circuit

import (
	"math/big"
	"testing"

	"github.com/akakou/zk-ban/snark"
	zkbanw "github.com/akakou/zk-ban/witness"
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

	_, err = gsk.IssueCertificate(upk)
	assert.NoError(err)

	commit, err := zkbanw.ComputeSignCommit(m, bsn, nonce, &usk)
	if err != nil {
		assert.NoError(err)
	}

	assign := &SyncCircuit{
		UserSecretKey:  usk.Number,
		Period:         period,
		RevocationList: [RevocationListSize]frontend.Variable{},
	}

	for i := 0; i < RevocationListSize; i++ {
		assign.RevocationList[i] = commit.Commit2
	}

	assign.GroupPublicKey.Assign(snark.TwistededwardsCurve, gpk.Bytes())
	// assign.Certificate.Assign(snark.TwistededwardsCurve, cert.Signature)

	assert.ProverSucceeded(&SyncCircuit{}, assign, test.WithCurves(snark.EcCurve))
}
