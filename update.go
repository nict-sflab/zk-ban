package zkban

import (
	"math/big"

	"github.com/akakou/zk-ban/circuit"
	"github.com/akakou/zk-ban/commit"
	"github.com/akakou/zk-ban/snark"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
)

func Update(nextPeriod *big.Int, signer *commit.Signer, rl commit.RevocationList, gpk *commit.GroupPublicKey, prover *snark.SnarkProver) (*commit.Signer, groth16.Proof, witness.Witness, error) {
	nextSigner, err := signer.Next(nextPeriod)
	if err != nil {
		return nil, nil, nil, err
	}

	w := circuit.NewUpdateCircuitWitness(nextSigner, signer, rl, gpk)

	proof, wit, _, err := snark.ProveSNARK(w, prover)
	if err != nil {
		return nil, nil, nil, err
	}

	return nextSigner, proof, wit, nil
}
