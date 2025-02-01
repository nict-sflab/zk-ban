package zkban

import (
	"crypto/rand"
	"math/big"

	"github.com/akakou/zk-ban/circuit"
	"github.com/akakou/zk-ban/commit"
	"github.com/akakou/zk-ban/snark"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
)

func JoinRequest(period *big.Int, snarkProver *snark.SnarkProver) (groth16.Proof, witness.Witness, *commit.UserSecretKey, *commit.UserPublicKey, error) {
	u, err := rand.Int(rand.Reader, max())
	if err != nil {
		return nil, nil, nil, nil, err
	}

	usk := commit.UserSecretKey{
		Number: u,
	}

	upk, err := usk.PublicKey(period)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	witness := circuit.NewJoinRequestWitness(period, usk.Number, upk.Number)

	proof, pubWit, _, err := snark.ProveSNARK(witness, snarkProver)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return proof, pubWit, &usk, upk, nil
}
