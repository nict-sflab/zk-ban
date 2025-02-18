package circuit

import (
	"math/big"
	"testing"

	"github.com/akakou/zk-ban/witness"
	"github.com/akakou/zk-ban/snark"
	"github.com/consensys/gnark/test"
)

func TestJoinReq(t *testing.T) {
	assert := test.NewAssert(t)
	var joinCercuit JoinRequestCircuit

	usk := witness.UserSecretKey{Number: big.NewInt(1)}
	period := big.NewInt(2024)

	upk, err := usk.PublicKey(period)
	assert.NoError(err)

	assign := NewJoinRequestWitness(
		period,
		upk.Number,
		usk.Number,
	)

	assert.ProverSucceeded(&joinCercuit, assign, test.WithCurves(snark.EcCurve))
}
