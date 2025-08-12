package precomputes

import (
	gnarkprecomputes "github.com/akakou/gnark-precomputes"
	zkban "github.com/akakou/zk-ban"
	"github.com/akakou/zk-ban/circuit"
	"github.com/akakou/zk-ban/primitives"
	zkbanw "github.com/akakou/zk-ban/witness"
	curve_bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381"
	fr_bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/twistededwards"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/twistededwards/eddsa"
	"github.com/consensys/gnark/backend/groth16"
	groth16_bls12381 "github.com/consensys/gnark/backend/groth16/bls12-381"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/frontend"
)

var UpdatePreparableIndex = 4

type UpdateCircuit struct {
	UpdateCircuit circuit.UpdateCircuit
}

func (circuit *UpdateCircuit) PreparableIndex() int {
	return UpdatePreparableIndex
}

func (circuit *UpdateCircuit) Define(api frontend.API) error {
	return circuit.UpdateCircuit.Define(api)
}

type PreparedUpdateRequestVerifyingKey[
	Vector any,
	G1Jac any,
	Proof any,
] struct {
	gnarkprecomputes.PreparedVerifyingKey[Vector, G1Jac, Proof]
}

func NewUpdateVerificationKeyBLS12381(gk groth16.VerifyingKey) (*PreparedUpdateRequestVerifyingKey[fr_bls12381.Vector, *curve_bls12381.G1Jac, *groth16_bls12381.Proof], error) {
	vk, err := gnarkprecomputes.FromBLS12381GnarkKey(gk, &UpdateCircuit{})
	if err != nil {
		return nil, err
	}

	return &PreparedUpdateRequestVerifyingKey[fr_bls12381.Vector, *curve_bls12381.G1Jac, *groth16_bls12381.Proof]{vk}, nil
}

func (vk *PreparedUpdateRequestVerifyingKey[Vector, G1Jac, Proof]) PrecomputeVerify(
	rl zkbanw.RevocationList,
	gpk *zkbanw.GroupPublicKey,
) (*G1Jac, error) {
	pubWit, err := newPublicPrecomputationUpdateCircuitWitness(rl, gpk)
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
) error {
	pubWit, err := newPublicPreparedUpdateCircuitWitness(nextPeriod, updateRequest.PublicKey, updateRequest.UpdateTicket, lastPeriod)
	if err != nil {
		return err
	}

	err = vk.PreparedVerifyingKey.VerifyPrepared(updateRequest.Proof.Proof, pubWit, prepare)
	return err
}

func newPublicPrecomputationUpdateCircuitWitness(
	revocationList zkbanw.RevocationList,
	gpk *zkbanw.GroupPublicKey,
) (witness.Witness, error) {
	wit, err := circuit.NewUpdateCircuitWitness(0,
		&zkbanw.UserPublicKey{primitives.NewBigInt(0)},
		&zkbanw.OneTimeTicket{primitives.NewBigInt(0)},
		&zkbanw.Signer{
			UserSecretKey: &zkbanw.UserSecretKey{
				primitives.NewBigInt(0),
			},
			Credential: &zkbanw.Credential{
				Signature: make([]byte, 32),
			},
			Period: 0,
		}, revocationList, gpk)

	if err != nil {
		return nil, err
	}

	return wit.Public()
}

func newPublicPreparedUpdateCircuitWitness(
	nextPeriod int64,
	nextPublicKey *zkbanw.UserPublicKey,
	updateTicket *zkbanw.OneTimeTicket,
	lastPeriod int64,
) (witness.Witness, error) {
	wit, err := circuit.NewUpdateCircuitWitness(nextPeriod, nextPublicKey, updateTicket, &zkbanw.Signer{
		UserSecretKey: &zkbanw.UserSecretKey{
			primitives.NewBigInt(0),
		},
		Credential: &zkbanw.Credential{
			Signature: make([]byte, 32),
		},
		Period: lastPeriod,
	}, zkbanw.EmptyRevocationList([]int{}), &zkbanw.GroupPublicKey{
		PublicKey: &eddsa.PublicKey{
			A: twistededwards.NewPointAffine([4]uint64{}, [4]uint64{}),
		},
	})

	if err != nil {
		return nil, err
	}

	return wit.Public()
}
