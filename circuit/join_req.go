package circuit

import (
	"math/big"

	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
)

type JoinRequestCircuit struct {
	UserSecretKey frontend.Variable `gnark:"sk"`
	UserPublicKey frontend.Variable `gnark:",public"`
	Period        frontend.Variable `gnark:",public"`
}

func (circuit *JoinRequestCircuit) Define(api frontend.API) error {
	mimc, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}

	mimc.Write(circuit.Period, circuit.UserSecretKey)
	h := mimc.Sum()

	api.AssertIsEqual(circuit.UserPublicKey, h)

	return nil
}

func JoinRequestWitness(usk *big.Int, upk []byte, period *big.Int) *JoinRequestCircuit {
	assign := &JoinRequestCircuit{
		UserSecretKey: usk,
		UserPublicKey: upk,
		Period:        period,
	}

	return assign
}
