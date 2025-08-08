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

type Proof struct{ groth16.Proof }
type ConstraintSystem struct{ constraint.ConstraintSystem }
type ProveKey struct{ groth16.ProvingKey }
type VerifyKey struct{ groth16.VerifyingKey }

type SnarkParams struct {
	ConstraintSystem ConstraintSystem
	ProveKey         ProveKey
	VerifyKey        VerifyKey
	Circuit          frontend.Circuit
}

func (params *SnarkParams) Prover() *SnarkProver {
	return &SnarkProver{
		ConstraintSystem: params.ConstraintSystem,
		ProveKey:         params.ProveKey,
	}
}

type SnarkProver struct {
	ConstraintSystem ConstraintSystem
	ProveKey         ProveKey
}
