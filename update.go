package zkban

import (
	"math/big"

	"github.com/akakou/zk-ban/circuit"
	"github.com/akakou/zk-ban/snark"
	zkbanw "github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/backend/groth16"
)

func UpdateRequest(nextPeriod *big.Int, signer *zkbanw.Signer, rl zkbanw.RevocationList, gpk *zkbanw.GroupPublicKey, prover *snark.SnarkProver) (*zkbanw.Signer, groth16.Proof, *circuit.UpdateCircuit, error) {
	nextSigner, err := signer.NextWithoutCred(nextPeriod, gpk)
	if err != nil {
		return nil, nil, nil, err
	}

	w := circuit.NewUpdateCircuitWitness(nextSigner, signer, rl)

	proof, _, _, err := snark.ProveSNARK(w, prover)
	if err != nil {
		return nil, nil, nil, err
	}

	return nextSigner, proof, w, nil
}
