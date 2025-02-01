package zkban

import (
	"math/big"

	"github.com/akakou/zk-ban/circuit"
	zkbanc "github.com/akakou/zk-ban/circuit"
	"github.com/akakou/zk-ban/snark"
	zkbanw "github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
)

func Sync(commit [zkbanc.SessionSize][circuit.RevocationPerSession]*big.Int, nextPeriod *big.Int, signer *zkbanw.Signer, sessionName [zkbanc.SessionSize]*big.Int, gpk *zkbanw.GroupPublicKey, prover *snark.SnarkProver) (groth16.Proof, witness.Witness, error) {
	nextSignerCandidate, err := signer.NextSignerCandidate(nextPeriod)
	if err != nil {
		return nil, nil, err
	}

	w := zkbanc.NewSyncCircuitWitness(commit, sessionName, nextSignerCandidate, signer, gpk)

	proof, wit, _, err := snark.ProveSNARK(w, prover)
	if err != nil {
		return nil, nil, err
	}

	return proof, wit, nil
}
