package snark

import (
	"math/big"
	"testing"

	zkbanw "github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark/test"
)

func TestProofCircuit(t *testing.T) {
	assert := test.NewAssert(t)

	usk := zkbanw.UserSecretKey{UserSecretKey: big.NewInt(1)}
	r := big.NewInt(2)
	m := big.NewInt(3)
	period := big.NewInt(2024)

	gsk, gpk, err := zkbanw.RandomGroupKeyPair()
	assert.NoError(err)

	upk, err := usk.PublicKey(period)
	assert.NoError(err)

	cert, err := gsk.IssuseCertificate(upk)
	assert.NoError(err)

	proof, err := zkbanw.ComputeProveWitness(m, r, &usk)
	assert.NoError(err)

	authCircuit := ProofCircuit{}

	assign := &ProofCircuit{
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
