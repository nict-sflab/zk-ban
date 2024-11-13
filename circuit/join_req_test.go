package circuit

import (
	"math/big"
	"testing"

	"github.com/akakou/zk-ban/snark"
	zkbanw "github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/test"
)

func TestJoinReq(t *testing.T) {
	assert := test.NewAssert(t)
	var joinCercuit JoinRequestCircuit

	usk := zkbanw.UserSecretKey{Number: big.NewInt(1)}
	period := big.NewInt(2024)

	upk, err := usk.PublicKey(period)
	assert.NoError(err)

	assign := &JoinRequestCircuit{
		UserSecretKey: usk.Number,
		UserPublicKey: upk.Buffer,
		Period:        period,
	}

	assert.ProverSucceeded(&joinCercuit, assign, test.WithCurves(snark.EcCurve))
}
