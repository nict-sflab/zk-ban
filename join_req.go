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
	UserPublicKey *zkbanw.UserPublicKey
	groth16.Proof
}

func RequestJoin(period int64, snarkProver *snark.SnarkProver) (*JoinRequest, *zkbanw.UserSecretKey, error) {
	u, err := rand.Int(rand.Reader, max())
	if err != nil {
		return nil, nil, err
	}

	usk := zkbanw.UserSecretKey{
		BigInt: &primitives.BigInt{*u},
	}

	upk, err := usk.PublicKey(period)
	if err != nil {
		return nil, nil, err
	}

	wit, err := circuit.NewJoinRequestWitness(period, upk, &usk)

	proof, err := groth16.Prove(snarkProver.ConstraintSystem.ConstraintSystem, snarkProver.ProveKey.ProvingKey, wit)
	if err != nil {
		return nil, nil, err
	}

	return &JoinRequest{UserPublicKey: upk, Proof: proof}, &usk, nil
}

func (req *JoinRequest) Verify(period int64, verifyKey snark.VerifyKey) error {
	wit, err := circuit.NewPublicJoinRequestWitness(period, req.UserPublicKey)
	if err != nil {
		return err
	}

	pubWit, err := wit.Public()
	if err != nil {
		return err
	}

	err = groth16.Verify(req.Proof, verifyKey.VerifyingKey, pubWit)
	return err
}
