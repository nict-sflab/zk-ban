package zkban

import (
	"github.com/akakou/zk-ban/circuit"
	"github.com/akakou/zk-ban/snark"
	zkbanw "github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/backend/groth16"
)

type UpdateRequest struct {
	PublicKey    *zkbanw.UserPublicKey
	UpdateTicket *zkbanw.OneTimeTicket
	Proof        groth16.Proof
}

func RequestUpdate(nextPeriod int64, signer *zkbanw.Signer, rl zkbanw.RevocationList, gpk *zkbanw.GroupPublicKey, prover *snark.SnarkProver) (*UpdateRequest, error) {
	ticket, err := signer.UserSecretKey.OneTimeTicket(signer.Period)
	if err != nil {
		return nil, err
	}

	nextPublicKey, err := signer.UserSecretKey.PublicKey(nextPeriod)
	if err != nil {
		return nil, err
	}

	wit, err := circuit.NewUpdateCircuitWitness(nextPeriod, nextPublicKey, ticket, signer, rl, gpk)
	if err != nil {
		return nil, err
	}

	proof, err := groth16.Prove(prover.ConstraintSystem, prover.ProveKey, wit)
	if err != nil {
		return nil, err
	}

	updateReq := UpdateRequest{
		Proof:        proof,
		PublicKey:    nextPublicKey,
		UpdateTicket: ticket,
	}

	return &updateReq, nil
}

func (request *UpdateRequest) Verify(nextPeriod int64, lastPeriod int64, rl zkbanw.RevocationList, gpk *zkbanw.GroupPublicKey, verifyKey groth16.VerifyingKey) error {
	pubWit, err := circuit.NewPublicUpdateCircuitWitness(nextPeriod, request.PublicKey, request.UpdateTicket, lastPeriod, rl, gpk)
	if err != nil {
		return err
	}

	err = groth16.Verify(request.Proof, verifyKey, pubWit)
	return err
}
