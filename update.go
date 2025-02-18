package zkban

import (
	"math/big"

	"github.com/akakou/zk-ban/circuit"
	"github.com/akakou/zk-ban/snark"
	zkbanw "github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
)

func Update(nextPeriod *big.Int, nextGpk *zkbanw.GroupPublicKey, signer *zkbanw.Signer, rl zkbanw.RevocationList, prover *snark.SnarkProver) (*zkbanw.Signer, groth16.Proof, witness.Witness, error) {
	nextSigner, err := signer.NextWithoutCred(nextPeriod, nextGpk)
	if err != nil {
		return nil, nil, nil, err
	}

	w := circuit.NewUpdateCircuitWitness(nextSigner, signer, rl)

	proof, wit, _, err := snark.ProveSNARK(w, prover)
	if err != nil {
		return nil, nil, nil, err
	}

	return nextSigner, proof, wit, nil
}
