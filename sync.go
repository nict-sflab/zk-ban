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

func Sync(basename *big.Int, commit [circuit.RevocationListSize]*big.Int, signer *zkbanw.Signer, gpk *zkbanw.GroupPublicKey, prover *snark.SnarkProver) (groth16.Proof, witness.Witness, error) {
	w := zkbanc.NewSyncCircuitWitness(commit, basename, signer, gpk)

	proof, wit, _, err := snark.ProveSNARK(w, prover)
	if err != nil {
		return nil, nil, err
	}

	return proof, wit, nil
}
