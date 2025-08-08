package zkban

import (
	"crypto/rand"

	"github.com/akakou/zk-ban/circuit"
	"github.com/akakou/zk-ban/primitives"
	"github.com/akakou/zk-ban/snark"
	zkbanw "github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/backend/groth16"
)

type JoinRequest struct {
	UserPublicKey *zkbanw.UserPublicKey `json:"user_public_key"`
	Proof         *snark.Proof          `json:"proof"`
}

func RequestJoin(period int64, snarkProver *snark.SnarkProver) (*JoinRequest, *zkbanw.UserSecretKey, error) {
	u, err := rand.Int(rand.Reader, max())
	if err != nil {
		return nil, nil, err
	}

	usk := zkbanw.UserSecretKey{
		Number: &primitives.BigInt{Int: *u},
	}

	upk, err := usk.PublicKey(period)
	if err != nil {
		return nil, nil, err
	}

	wit, err := circuit.NewJoinRequestWitness(period, upk, &usk)
	if err != nil {
		return nil, nil, err
	}

	proof, err := groth16.Prove(snarkProver.ConstraintSystem, snarkProver.ProveKey, wit)
	if err != nil {
		return nil, nil, err
	}

	return &JoinRequest{UserPublicKey: upk, Proof: &snark.Proof{proof}}, &usk, nil
}

func (req *JoinRequest) Verify(period int64, verifyKey groth16.VerifyingKey) error {
	wit, err := circuit.NewPublicJoinRequestWitness(period, req.UserPublicKey)
	if err != nil {
		return err
	}

	pubWit, err := wit.Public()
	if err != nil {
		return err
	}

	err = groth16.Verify(req.Proof.Proof, verifyKey, pubWit)
	return err
}
