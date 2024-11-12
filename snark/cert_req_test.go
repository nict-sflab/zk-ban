package snark

import (
	"math/big"
	"testing"

	zkban "github.com/akakou/zk-ban"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/test"
)

func TestJoin(t *testing.T) {
	assert := test.NewAssert(t)
	var joinCercuit CertificateRequestCircuit

	usk := zkban.UserSecretKey{UserSecretKey: big.NewInt(1)}
	period := big.NewInt(2024)

	upk, err := usk.PublicKey(period)
	assert.NoError(err)

	assign := &CertificateRequestCircuit{
		UserSecretKey: usk.UserSecretKey,
		UserPublicKey: upk.UserPublicKey,
		Period:        period,
	}

	assert.ProverSucceeded(&joinCercuit, assign, test.WithCurves(ecc.BN254))
}
