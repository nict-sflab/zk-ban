package zkban

import (
	"crypto/rand"
	"math/big"

	zkbanc "github.com/akakou/zk-ban/circuit"
	zkbanw "github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
)

func max() *big.Int {
	var i, e = big.NewInt(2), big.NewInt(32)
	i.Exp(i, e, nil)

	return i
}

func Prove(m *big.Int, usk *zkbanw.UserSecretKey, upk *zkbanw.UserPublicKey, cert *zkbanw.Certificate, gpk *zkbanw.GroupPublicKey, pk groth16.ProvingKey, ccs constraint.ConstraintSystem) (groth16.Proof, *zkbanw.Proof, error) {
	r, err := rand.Int(rand.Reader, max())
	if err != nil {
		return nil, nil, err
	}

	h, err := zkbanw.ComputeProveWitness(m, r, usk)
	if err != nil {
		return nil, nil, err
	}

	w := zkbanc.ProofWitness(m, r, usk.UserSecretKey, upk.UserPublicKey, h.Hash, cert.Signature, gpk)

	proof, err := proveSNARK(w, pk, ccs)
	if err != nil {
		return nil, nil, err
	}

	return proof, h, nil
}
