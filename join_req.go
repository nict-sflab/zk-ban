package zkban

import (
	"crypto/rand"
	"math/big"

	"github.com/akakou/zk-ban/circuit"
	"github.com/akakou/zk-ban/snark"
	zkbanw "github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/backend/groth16"
)

func JoinRequest(period *big.Int, snarkProver *snark.SnarkProver) (groth16.Proof, *circuit.JoinRequestCircuit, error) {
	u, err := rand.Int(rand.Reader, max())
	if err != nil {
		return nil, nil, err
	}

	usk := zkbanw.UserSecretKey{
		Number: u,
	}

	upk, err := usk.PublicKey(period)
	if err != nil {
		return nil, nil, err
	}

	assign := circuit.NewJoinRequestWitness(period, upk.Number, usk.Number)

	proof, _, _, err := snark.ProveSNARK(assign, snarkProver)
	if err != nil {
		return nil, nil, err
	}

	return proof, assign, nil
}
