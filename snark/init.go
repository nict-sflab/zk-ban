package snark

import (
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

func InitSNARK[T frontend.Circuit](w T) (*SnarkParams, error) {
	ccs, err := frontend.Compile(EcCurve.ScalarField(), r1cs.NewBuilder, w)

	if err != nil {
		return nil, err
	}

	pk, vk, err := groth16.Setup(ccs)

	if err != nil {
		return nil, err
	}

	return &SnarkParams{ccs, pk, vk, w}, nil
}
