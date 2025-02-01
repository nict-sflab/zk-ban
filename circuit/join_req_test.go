package circuit

import (
	"math/big"
	"testing"

	"github.com/akakou/zk-ban/commit"
	"github.com/akakou/zk-ban/snark"
	"github.com/consensys/gnark/test"
)

func TestJoinReq(t *testing.T) {
	assert := test.NewAssert(t)
	var joinCercuit JoinRequestCircuit

	usk := commit.UserSecretKey{Number: big.NewInt(1)}
	period := big.NewInt(2024)

	upk, err := usk.PublicKey(period)
	assert.NoError(err)

	assign := &JoinRequestCircuit{
		UserSecretKey: usk.Number,
		UserPublicKey: upk.Number,
		Period:        period,
	}

	assert.ProverSucceeded(&joinCercuit, assign, test.WithCurves(snark.EcCurve))
}
