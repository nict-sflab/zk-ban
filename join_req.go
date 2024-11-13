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

func JoinRequest(period *big.Int, snarkProver *snark.SnarkProver) (groth16.Proof, witness.Witness, *zkbanw.UserSecretKey, *zkbanw.UserPublicKey, error) {
	u, err := rand.Int(rand.Reader, max())
	if err != nil {
		return nil, nil, nil, nil, err
	}

	usk := zkbanw.UserSecretKey{
		Number: u,
	}

	upk, err := usk.PublicKey(period)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	witness := zkbanc.JoinRequestWitness(period, usk.Number, upk.Buffer)

	proof, pubWit, _, err := snark.ProveSNARK(witness, snarkProver)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return proof, pubWit, &usk, upk, nil
}
