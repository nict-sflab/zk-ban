package circuit

import (
	"math/big"

	"github.com/consensys/gnark/frontend"
)

type JoinRequestCircuit struct {
	UserSecretKey     frontend.Variable `gnark:",secret"`
	PublicKeyAuthInfo PublicKeyAuthInfo
}

func (circuit *JoinRequestCircuit) Define(api frontend.API) error {
	err := authPubKey(api, circuit.PublicKeyAuthInfo, circuit.UserSecretKey)

	if err != nil {
		return err
	}
	return nil
}

func NewJoinRequestWitness(period, upk, usk *big.Int) *JoinRequestCircuit {
	assign := &JoinRequestCircuit{
		UserSecretKey: usk,
		PublicKeyAuthInfo: PublicKeyAuthInfo{
			UserPublicKey: upk,
			Period:        period,
		},
	}

	return assign
}
