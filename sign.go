package zkban

import (
	"math/big"

	"github.com/akakou/zk-ban/circuit"
	"github.com/akakou/zk-ban/snark"
	zkbanw "github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/backend/groth16"
)

func max() *big.Int {
	var i, e = big.NewInt(2), big.NewInt(32)
	i.Exp(i, e, nil)

	return i
}

func Sign(m, counter *big.Int, signer *zkbanw.Signer, prover *snark.SnarkProver) (groth16.Proof, *circuit.SignCircuit, error) {
	sessionTag := signer.SessionTag(counter)

	comm, err := signer.CommitSign(m, sessionTag)
	if err != nil {
		return nil, nil, err
	}

	w := circuit.NewSignWitness(m, sessionTag, comm, signer)

	proof, _, _, err := snark.ProveSNARK(w, prover)
	if err != nil {
		return nil, nil, err
	}

	return proof, w, nil
}
