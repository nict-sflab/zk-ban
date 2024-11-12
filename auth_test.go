package zkban

import (
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark-crypto/hash"
	"github.com/consensys/gnark-crypto/signature"
	"github.com/consensys/gnark-crypto/signature/eddsa"
	"github.com/consensys/gnark/test"
)

func setup(assert *test.Assert) (signature.Signer, signature.PublicKey) {
	gsk, err := eddsa.New(twistededwards.BN254, rand.Reader)
	assert.NoError(err)
	gpk := gsk.Public()

	return gsk, gpk
}

func join(sk *big.Int, gsk signature.Signer, assert *test.Assert) ([]byte, []byte) {
	hFunc := hash.MIMC_BN254.New()
	period := big.NewInt(2024)

	hasher := mimc.NewMiMC(mimc.WithByteOrder(fr.BigEndian))
	_, err := hasher.Write(sk.Bytes())
	assert.NoError(err)

	_, err = hasher.Write(period.Bytes())
	assert.NoError(err)

	pk := hasher.Sum(nil)

	cert, err := gsk.Sign(pk, hFunc)
	assert.NoError(err)

	return pk, cert
}

func prove(m, r, sk *big.Int, assert *test.Assert) []byte {
	hasher := mimc.NewMiMC(mimc.WithByteOrder(fr.BigEndian))
	_, err := hasher.Write(m.Bytes())
	assert.NoError(err)

	_, err = hasher.Write(r.Bytes())
	assert.NoError(err)

	_, err = hasher.Write(sk.Bytes())
	assert.NoError(err)

	h := hasher.Sum(nil)

	return h
}

func TestAuth(t *testing.T) {
	assert := test.NewAssert(t)

	sk := big.NewInt(1)
	r := big.NewInt(2)
	m := big.NewInt(3)

	gsk, gpk := setup(assert)
	upk, cert := join(sk, gsk, assert)
	hash := prove(m, r, sk, assert)

	authCircuit := AuthCircuit{}

	assign := &AuthCircuit{
		Nonce:         r,
		UserSecretKey: sk,
		Message:       m,
		Hash:          hash,
		UserPublicKey: upk,
	}

	assign.GroupPublicKey.Assign(twistededwards.BN254, gpk.Bytes())
	assign.Certificate.Assign(twistededwards.BN254, cert)

	assert.ProverSucceeded(&authCircuit, assign, test.WithCurves(ecc.BN254))
}
