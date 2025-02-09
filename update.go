package zkban

import (
	"math/big"

	"github.com/akakou/zk-ban/circuit"
	"github.com/akakou/zk-ban/commit"
	"github.com/akakou/zk-ban/snark"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
)

func Update(nextPeriod *big.Int, nextGpk *commit.GroupPublicKey, signer *commit.Signer, rl commit.RevocationList, prover *snark.SnarkProver) (*commit.Signer, groth16.Proof, witness.Witness, error) {
	nextSigner, err := signer.NextWithoutCred(nextPeriod, nextGpk)
	if err != nil {
		return nil, nil, nil, err
	}

	w := circuit.NewUpdateCircuitWitness(nextSigner, signer, rl, signer.GroupPublicKey)

	proof, wit, _, err := snark.ProveSNARK(w, prover)
	if err != nil {
		return nil, nil, nil, err
	}

	return nextSigner, proof, wit, nil
}
