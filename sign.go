package zkban

import (
	"crypto/rand"
	"math/big"

	zkbanc "github.com/akakou/zk-ban/circuit"
	"github.com/akakou/zk-ban/snark"
	zkbanw "github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
)

func max() *big.Int {
	var i, e = big.NewInt(2), big.NewInt(32)
	i.Exp(i, e, nil)

	return i
}

func Sign(m *big.Int, signer *zkbanw.Signer, gpk *zkbanw.GroupPublicKey, prover *snark.SnarkProver) (groth16.Proof, witness.Witness, error) {
	r, err := rand.Int(rand.Reader, max())
	if err != nil {
		return nil, nil, err
	}

	h, err := zkbanw.ComputeSignCommit(m, r, signer.UserSecretKey)
	if err != nil {
		return nil, nil, err
	}

	w := zkbanc.NewSignWitness(m, r, h.Buffer, signer, gpk)

	proof, wit, _, err := snark.ProveSNARK(w, prover)
	if err != nil {
		return nil, nil, err
	}

	return proof, wit, nil
}
