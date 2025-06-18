package zkban

import (
	"fmt"
	"math/big"

	"github.com/akakou/zk-ban/circuit"
	"github.com/akakou/zk-ban/snark"
	zkbanw "github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/backend/groth16"
)

func UpdateRequest(nextPeriod *big.Int, signer *zkbanw.Signer, rl zkbanw.RevocationList, gpk *zkbanw.GroupPublicKey, prover *snark.SnarkProver) (*zkbanw.Signer, groth16.Proof, *circuit.UpdateCircuit, error) {
	nextSigner, err := signer.NextWithoutCred(nextPeriod)
	if err != nil {
		return nil, nil, nil, err
	}

	rlHash := zkbanw.HashRevocationList(rl)
	fmt.Printf("rlHash%v\n", rlHash.Text(10))

	w := circuit.NewUpdateCircuitWitness(nextSigner, signer, rl, gpk)
	w.RevocationListHash = rlHash

	proof, _, _, err := snark.ProveSNARK(w, prover)
	if err != nil {
		return nil, nil, nil, err
	}

	return nextSigner, proof, w, nil
}
