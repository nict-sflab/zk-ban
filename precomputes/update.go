package precomputes

import (
	gnarkprecomputes "github.com/akakou/gnark-precomputes"
	zkban "github.com/akakou/zk-ban"
	"github.com/akakou/zk-ban/circuit"
	zkbanw "github.com/akakou/zk-ban/witness"
	curve_bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381"
	fr_bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
	"github.com/consensys/gnark/backend/groth16"
	groth16_bls12381 "github.com/consensys/gnark/backend/groth16/bls12-381"
)

type PreparedUpdateRequestVerifyingKey[
	Vector any,
	G1Jac any,
	Proof any,
] struct {
	gnarkprecomputes.PreparedVerifyingKey[Vector, G1Jac, Proof]
}

func NewUpdateVerificationKeyBLS12381(gk groth16.VerifyingKey) (*PreparedUpdateRequestVerifyingKey[fr_bls12381.Vector, *curve_bls12381.G1Jac, *groth16_bls12381.Proof], error) {
	vk, err := gnarkprecomputes.FromBLS12381GnarkKey(gk, &circuit.UpdateCircuit{})
	if err != nil {
		return nil, err
	}

	return &PreparedUpdateRequestVerifyingKey[fr_bls12381.Vector, *curve_bls12381.G1Jac, *groth16_bls12381.Proof]{vk}, nil
}

func (vk *PreparedUpdateRequestVerifyingKey[Vector, G1Jac, Proof]) PrecomputeVerifyingUpdateRequest(
	updateRequest *zkban.UpdateRequest,
	nextPeriod,
	lastPeriod int64,
	rl zkbanw.RevocationList,
	gpk *zkbanw.GroupPublicKey,
) (*G1Jac, error) {
	pubWit, err := circuit.NewPublicUpdateCircuitWitness(nextPeriod, updateRequest.PublicKey, updateRequest.UpdateTicket, lastPeriod, rl, gpk)
	if err != nil {
		return nil, err
	}

	prepare, err := vk.PreparePublicInputs(pubWit)
	if err != nil {
		return nil, err
	}

	return &prepare, nil
}

func (vk *PreparedUpdateRequestVerifyingKey[Vector, G1Jac, Proof]) VerifyPrepared(
	prepare G1Jac,
	updateRequest *zkban.UpdateRequest,
	nextPeriod,
	lastPeriod int64,
	rl zkbanw.RevocationList,
	gsk *zkbanw.GroupSecretKey,
	gpk *zkbanw.GroupPublicKey,
) error {
	pubWit, err := circuit.NewPublicUpdateCircuitWitness(nextPeriod, updateRequest.PublicKey, updateRequest.UpdateTicket, lastPeriod, rl, gpk)
	if err != nil {
		return err
	}

	err = vk.PreparedVerifyingKey.VerifyPrepared(updateRequest.Proof, pubWit, prepare)
	return err
}
