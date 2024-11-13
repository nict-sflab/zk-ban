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

func TestAll(t *testing.T) {
	assert := test.NewAssert(t)
	gsk, gpk, err := zkbanw.RandomGroupKeyPair()
	assert.NoError(err)

	ccs, pk, vk, err := snark.InitSNARK(&zkbanc.JoinRequestCircuit{})
	assert.NoError(err)

	period := big.NewInt(2024)

	var usk *zkbanw.UserSecretKey
	var upk *zkbanw.UserPublicKey

	t.Run("join req", func(t *testing.T) {
		var proof groth16.Proof
		proof, pubWit, _usk, _upk, err := JoinRequest(period, pk, ccs)
		assert.NoError(err)

		err = groth16.Verify(proof, vk, pubWit)
		assert.NoError(err)

		usk = _usk
		upk = _upk
	})

	ccs, pk, vk, err = snark.InitSNARK(&zkbanc.ProofCircuit{})
	assert.NoError(err)

	cert, err := gsk.IssueCertificate(upk)
	assert.NoError(err)

	m := big.NewInt(100)

	t.Run("sign", func(t *testing.T) {
		proof, pubWit, err := Sign(m, usk, upk, cert, gpk, pk, ccs)
		assert.NoError(err)

		err = groth16.Verify(proof, vk, pubWit)
		assert.NoError(err)
	})
}
