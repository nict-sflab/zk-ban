package snark

import (
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
)

var EcCurve = ecc.BN254
var TwistededwardsCurve = twistededwards.BN254

type SnarkParams struct {
	ConstraintSystem constraint.ConstraintSystem
	ProveKey         groth16.ProvingKey
	VerifyKey        groth16.VerifyingKey
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

// type SnarkVerifier struct {
// 	ccs constraint.ConstraintSystem
// 	vk  groth16.VerifyingKey
// }
