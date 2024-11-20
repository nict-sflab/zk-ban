package zkban

import (
	"math/big"
	"testing"

	zkbanc "github.com/akakou/zk-ban/circuit"
	"github.com/akakou/zk-ban/snark"
	zkbanw "github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/test"
)

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

type TestParams struct {
	gpk       *zkbanw.GroupPublicKey
	gsk       *zkbanw.GroupSecretKey
	joinSnark *snark.SnarkParams
	signSnark *snark.SnarkParams
	upk       *zkbanw.UserPublicKey
	usk       *zkbanw.UserSecretKey
	m         *big.Int
	period    *big.Int
	bsn       *big.Int
	cert      *zkbanw.Certificate
}

func (params *TestParams) signer() *zkbanw.Signer {
	signer := zkbanw.Signer{
		UserSecretKey: params.usk,
		UserPublicKey: params.upk,
		Certificate:   params.cert,
		Period:        params.period,
	}

	return &signer
}

func prepare() TestParams {
	gsk, gpk, err := zkbanw.RandomGroupKeyPair()
	panicIfErr(err)

	joinSnark, err := snark.InitSNARK(&zkbanc.JoinRequestCircuit{})
	panicIfErr(err)

	signSnark, err := snark.InitSNARK(&zkbanc.SignCircuit{})
	panicIfErr(err)

	period := big.NewInt(2024)

	var usk = &zkbanw.UserSecretKey{
		Number: big.NewInt(100),
	}

	upk, err := usk.PublicKey(period)
	panicIfErr(err)

	m := big.NewInt(100)
	bsn := big.NewInt(100)

	cert, err := gsk.IssueCertificate(upk)
	panicIfErr(err)

	return TestParams{
		gpk:       gpk,
		gsk:       gsk,
		usk:       usk,
		upk:       upk,
		joinSnark: joinSnark,
		signSnark: signSnark,
		m:         m,
		bsn:       bsn,
		period:    period,
		cert:      cert,
	}
}

func TestAll(t *testing.T) {
	assert := test.NewAssert(t)
	params := prepare()

	t.Run("join req", func(t *testing.T) {
		var proof groth16.Proof
		proof, pubWit, _, _, err := JoinRequest(params.period, params.joinSnark.Prover())
		assert.NoError(err)

		err = groth16.Verify(proof, params.joinSnark.VerifyKey, pubWit)
		assert.NoError(err)
	})

	t.Run("sign", func(t *testing.T) {
		proof, pubWit, err := Sign(params.m, params.bsn, params.signer(), params.gpk, params.signSnark.Prover())
		assert.NoError(err)

		err = groth16.Verify(proof, params.signSnark.VerifyKey, pubWit)
		assert.NoError(err)
	})
}

func BenchmarkAll(t *testing.B) {
	// assert := test.NewAssert(t.(*testing.T))
	params := prepare()

	var proof groth16.Proof
	var pubWit witness.Witness
	var err error

	t.Run("join req", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			proof, pubWit, _, _, err = JoinRequest(params.period, params.joinSnark.Prover())
			panicIfErr(err)
		}
	})

	t.Run("verify join req", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			err = groth16.Verify(proof, params.joinSnark.VerifyKey, pubWit)
			panicIfErr(err)
		}
	})

	t.Run("sign", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			proof, pubWit, err = Sign(params.m, params.bsn, params.signer(), params.gpk, params.signSnark.Prover())
			panicIfErr(err)
		}
	})

	t.Run("verify", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			err = groth16.Verify(proof, params.signSnark.VerifyKey, pubWit)
			panicIfErr(err)
		}
	})
}
