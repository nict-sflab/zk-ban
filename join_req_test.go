package zkban

import (
	"math/big"
	"testing"

	zkbanc "github.com/akakou/zk-ban/circuit"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/test"
)

func TestJoinRequest(t *testing.T) {
	assert := test.NewAssert(t)

	joinReqCircuit := zkbanc.JoinRequestCircuit{}
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &joinReqCircuit)

	assert.NoError(err)

	pk, vk, err := groth16.Setup(ccs)
	assert.NoError(err)

	period := big.NewInt(2024)
	proof, pubWit, _, err := JoinRequest(period, pk, ccs)
	assert.NoError(err)

	err = groth16.Verify(proof, vk, pubWit)
	assert.NoError(err)
}
