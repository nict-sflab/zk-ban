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

type AuthCircuit struct {
	AuthCircuit circuit.SignCircuit
}

func (circuit *AuthCircuit) NonPrecomputables() []int {
	return []int{0, 1, 2}
}

func (circuit *AuthCircuit) Define(api frontend.API) error {
	return circuit.AuthCircuit.Define(api)
}

type PreparedAuthVerifyingKey[
	Vector any,
	G1Jac any,
	Proof any,
] struct {
	gnarkprecomputes.PreparedVerifyingKey[Vector, G1Jac, Proof]
}

func NewAuthVerificationKeyBLS12381(gk groth16.VerifyingKey) (*PreparedAuthVerifyingKey[fr_bls12381.Vector, *curve_bls12381.G1Jac, *groth16_bls12381.Proof], error) {
	vk, err := gnarkprecomputes.FromBLS12381GnarkKey(gk, &AuthCircuit{})
	if err != nil {
		return nil, err
	}

	return &PreparedAuthVerifyingKey[fr_bls12381.Vector, *curve_bls12381.G1Jac, *groth16_bls12381.Proof]{vk}, nil
}

func (vk *PreparedAuthVerifyingKey[Vector, G1Jac, Proof]) PrecomputeVerify(
	period int64,
	gpk *zkbanw.GroupPublicKey,
) (*G1Jac, error) {
	pubWit, err := newPublicPrecomputationAuthCircuitWitness(period, gpk)
	if err != nil {
		return nil, err
	}

	prepare, err := vk.PreparePublicInputs(pubWit)
	if err != nil {
		return nil, err
	}

	return &prepare, nil
}

func (vk *PreparedAuthVerifyingKey[Vector, G1Jac, Proof]) VerifyPrepared(
	prepare G1Jac,
	message *primitives.BigInt,
	authRequest *zkban.Signature,
) error {
	pubWit, err := newPublicPreparedAuthCircuitWitness(message, authRequest.Commit)
	if err != nil {
		return err
	}

	err = vk.PreparedVerifyingKey.VerifyPrepared(authRequest.Proof, pubWit, prepare)
	return err
}

func newPublicPrecomputationAuthCircuitWitness(
	currentPeriod int64,
	gpk *zkbanw.GroupPublicKey,
) (witness.Witness, error) {
	wit, err := circuit.NewSignWitness(
		primitives.NewBigInt(int64(0)),
		primitives.NewBigInt(int64(0)),
		0,
		&zkbanw.SignCommit{
			Sigma: primitives.NewBigInt(int64(0)),
			Nym:   primitives.NewBigInt(int64(0)),
		},
		&zkbanw.Signer{
			UserSecretKey: &zkbanw.UserSecretKey{
				Number: primitives.NewBigInt(int64(0)),
			},
			Credential: &zkbanw.Credential{
				Signature: make([]byte, 32),
			},
			Period: currentPeriod,
		},
		gpk)

	if err != nil {
		return nil, err
	}

	return wit.Public()
}

func newPublicPreparedAuthCircuitWitness(m *primitives.BigInt, commit *zkbanw.SignCommit) (witness.Witness, error) {
	wit, err := circuit.NewSignWitness(m, primitives.NewBigInt(0), 0, commit,
		&zkbanw.Signer{
			UserSecretKey: &zkbanw.UserSecretKey{
				primitives.NewBigInt(0),
			},
			Credential: &zkbanw.Credential{
				Signature: make([]byte, 32),
			},
			Period: 0,
		},
		&zkbanw.GroupPublicKey{
			PublicKey: &eddsa.PublicKey{
				A: twistededwards.NewPointAffine([4]uint64{}, [4]uint64{}),
			},
		})

	if err != nil {
		return nil, err
	}

	return wit.Public()
}
