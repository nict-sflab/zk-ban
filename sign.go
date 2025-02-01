package zkban

import (
	"math/big"

	"github.com/akakou/zk-ban/circuit"
	"github.com/akakou/zk-ban/commit"
	"github.com/akakou/zk-ban/snark"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
)

func max() *big.Int {
	var i, e = big.NewInt(2), big.NewInt(32)
	i.Exp(i, e, nil)

	return i
}

func Sign(m, bsn *big.Int, signer *commit.Signer, gpk *commit.GroupPublicKey, prover *snark.SnarkProver) (groth16.Proof, witness.Witness, error) {
	comm, err := signer.ComputeSignCommit(m, bsn)
	if err != nil {
		return nil, nil, err
	}

	w := circuit.NewSignWitness(m, bsn, comm, signer, gpk)

	proof, wit, _, err := snark.ProveSNARK(w, prover)
	if err != nil {
		return nil, nil, err
	}

	return proof, wit, nil
}
