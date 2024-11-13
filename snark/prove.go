package snark

import (
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/frontend"
)

func ProveSNARK[T frontend.Circuit](w T, prover *SnarkProver) (groth16.Proof, witness.Witness, witness.Witness, error) {
	witness, err := frontend.NewWitness(w, EcCurve.ScalarField())
	if err != nil {
		return nil, nil, nil, err
	}

	publicWit, err := witness.Public()
	if err != nil {
		return nil, nil, nil, err
	}

	proof, err := groth16.Prove(prover.ConstraintSystem, prover.ProveKey, witness)
	if err != nil {
		return nil, nil, nil, err
	}

	return proof, publicWit, witness, nil
}
