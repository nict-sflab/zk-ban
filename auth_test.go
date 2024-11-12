package zkban

import (
	"math/big"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark/test"
)

func TestAuth(t *testing.T) {
	assert := test.NewAssert(t)
	authCircuit := AuthCircuit{}

	sk := big.NewInt(1)
	r := big.NewInt(2)
	m := big.NewInt(3)

	hasher := mimc.NewMiMC(mimc.WithByteOrder(fr.BigEndian))
	_, err := hasher.Write(m.Bytes())
	assert.NoError(err)

	_, err = hasher.Write(r.Bytes())
	assert.NoError(err)

	_, err = hasher.Write(sk.Bytes())
	assert.NoError(err)

	h := hasher.Sum(nil)

	assert.ProverSucceeded(&authCircuit, &AuthCircuit{
		Nonce:     r,
		SecretKey: sk,
		Message:   m,
		Hash:      h,
	}, test.WithCurves(ecc.BN254))
}
