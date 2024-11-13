package zkban

import (
	"math/big"
	"testing"

	zkbanc "github.com/akakou/zk-ban/circuit"
	"github.com/akakou/zk-ban/snark"
	zkbanw "github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/test"
)

type TestParams struct {
	gpk       *zkbanw.GroupPublicKey
	gsk       *zkbanw.GroupSecretKey
	joinSnark *snark.SnarkParams
	signSnark *snark.SnarkParams
	upk       *zkbanw.UserPublicKey
	usk       *zkbanw.UserSecretKey
	m         *big.Int
	period    *big.Int
	cert      *zkbanw.Certificate
}

func (params *TestParams) signer() *zkbanw.Signer {
	signer := zkbanw.Signer{
		UserSecretKey: params.usk,
		UserPublicKey: params.upk,
		Certificate:   params.cert,
	}

	return &signer
}

func prepare(assert *test.Assert) TestParams {
	gsk, gpk, err := zkbanw.RandomGroupKeyPair()
	assert.NoError(err)

	joinSnark, err := snark.InitSNARK(&zkbanc.JoinRequestCircuit{})
	assert.NoError(err)

	signSnark, err := snark.InitSNARK(&zkbanc.SignCircuit{})
	assert.NoError(err)

	period := big.NewInt(2024)

	var usk = &zkbanw.UserSecretKey{
		Number: big.NewInt(100),
	}

	upk, err := usk.PublicKey(period)
	assert.NoError(err)

	m := big.NewInt(100)

	cert, err := gsk.IssueCertificate(upk)
	assert.NoError(err)

	return TestParams{
		gpk:       gpk,
		gsk:       gsk,
		usk:       usk,
		upk:       upk,
		joinSnark: joinSnark,
		signSnark: signSnark,
		m:         m,
		period:    period,
		cert:      cert,
	}
}

func TestAll(t *testing.T) {
	assert := test.NewAssert(t)
	params := prepare(assert)

	t.Run("join req", func(t *testing.T) {
		var proof groth16.Proof
		proof, pubWit, _, _, err := JoinRequest(params.period, params.joinSnark.Prover())
		assert.NoError(err)

		err = groth16.Verify(proof, params.joinSnark.VerifyKey, pubWit)
		assert.NoError(err)
	})

	t.Run("sign", func(t *testing.T) {
		proof, pubWit, err := Sign(params.m, params.signer(), params.gpk, params.signSnark.Prover())
		assert.NoError(err)

		err = groth16.Verify(proof, params.signSnark.VerifyKey, pubWit)
		assert.NoError(err)
	})
}
