package snark

import (
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
)

var EcCurve = ecc.BLS12_381
var TwistededwardsCurve = twistededwards.BLS12_381

type SnarkParams struct {
	ConstraintSystem constraint.ConstraintSystem
	ProveKey         groth16.ProvingKey
	VerifyKey        groth16.VerifyingKey
	Circuit          frontend.Circuit
}

func (params *SnarkParams) Prover() *SnarkProver {
	return &SnarkProver{
		ConstraintSystem: params.ConstraintSystem,
		ProveKey:         params.ProveKey,
	}
}

type SnarkProver struct {
	ConstraintSystem constraint.ConstraintSystem
	ProveKey         groth16.ProvingKey
}
