package zkban

import (
	"crypto/rand"
	"math/big"

	zkbanc "github.com/akakou/zk-ban/circuit"
	zkbanw "github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/constraint"
)

func JoinRequest(period *big.Int, pk groth16.ProvingKey, ccs constraint.ConstraintSystem) (groth16.Proof, witness.Witness, *zkbanw.UserSecretKey, *zkbanw.UserPublicKey, error) {
	u, err := rand.Int(rand.Reader, max())
	if err != nil {
		return nil, nil, nil, nil, err
	}

	usk := zkbanw.UserSecretKey{
		UserSecretKey: u,
	}

	upk, err := usk.PublicKey(period)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	witness := zkbanc.JoinRequestWitness(usk.UserSecretKey, upk.UserPublicKey, period)

	proof, pubWit, _, err := proveSNARK(witness, pk, ccs)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return proof, pubWit, &usk, upk, nil
}
