package circuit

import (
	"math/big"

	"github.com/consensys/gnark/frontend"
)

type JoinRequestCircuit struct {
	UserSecretKey frontend.Variable `gnark:"sk"`
	UserPublicKey frontend.Variable `gnark:",public"`
	Period        frontend.Variable `gnark:",public"`
}

func (circuit *JoinRequestCircuit) Define(api frontend.API) error {
	err := pubKeyAuth(api, circuit.Period, circuit.UserSecretKey, circuit.UserPublicKey)

	if err != nil {
		return err
	}
	return nil
}

func NewJoinRequestWitness(period, usk, upk *big.Int) *JoinRequestCircuit {
	assign := &JoinRequestCircuit{
		UserSecretKey: usk,
		UserPublicKey: upk,
		Period:        period,
	}

	return assign
}
