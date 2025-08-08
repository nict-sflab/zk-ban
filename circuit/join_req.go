package circuit

import (
	"github.com/akakou/zk-ban/primitives"
	"github.com/akakou/zk-ban/snark"
	zkbanw "github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/frontend"
)

type JoinRequestCircuit struct {
	UserSecretKey     frontend.Variable `gnark:",secret"`
	PublicKeyAuthInfo PublicKeyAuthInfo
}

func (circuit *JoinRequestCircuit) Define(api frontend.API) error {
	err := authPubKey(api, circuit.PublicKeyAuthInfo, circuit.UserSecretKey, PUBLIC_KEY)

	if err != nil {
		return err
	}
	return nil
}

func NewJoinRequestWitness(period int64, upk *zkbanw.UserPublicKey, usk *zkbanw.UserSecretKey) (witness.Witness, error) {
	assign := &JoinRequestCircuit{
		UserSecretKey: usk.Int,
		PublicKeyAuthInfo: PublicKeyAuthInfo{
			UserPublicKey: upk.Int,
			Period:        period,
		},
	}

	return frontend.NewWitness(assign, snark.EcCurve.ScalarField())
}

func NewPublicJoinRequestWitness(period int64, upk *zkbanw.UserPublicKey) (witness.Witness, error) {
	wit, err := NewJoinRequestWitness(period, upk, &zkbanw.UserSecretKey{
		BigInt: primitives.NewBigInt(0),
	})

	if err != nil {
		return nil, err
	}

	return wit.Public()
}
