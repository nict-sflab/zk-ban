package zkban

import (
	"math/big"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark/test"
)

func TestAuth(t *testing.T) {
	assert := test.NewAssert(t)

	usk := UserSecretKey{big.NewInt(1)}
	r := big.NewInt(2)
	m := big.NewInt(3)
	period := big.NewInt(2024)

	gsk, gpk, err := RandomGroupKeyPair()
	assert.NoError(err)

	upk, err := usk.PublicKey(period)
	assert.NoError(err)

	cert, err := gsk.IssuseCertificate(upk)
	assert.NoError(err)

	proof, err := Prove(m, r, &usk)
	assert.NoError(err)

	authCircuit := AuthCircuit{}

	assign := &AuthCircuit{
		Nonce:         r,
		UserSecretKey: usk.UserSecretKey,
		Message:       m,
		Hash:          proof.Hash,
		UserPublicKey: upk.UserPublicKey,
	}

	assign.GroupPublicKey.Assign(twistededwards.BN254, gpk.Bytes())
	assign.Certificate.Assign(twistededwards.BN254, cert.Signature)

	assert.ProverSucceeded(&authCircuit, assign, test.WithCurves(ecc.BN254))
}
