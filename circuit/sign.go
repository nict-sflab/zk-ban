package circuit

import (
	"math/big"

	"github.com/consensys/gnark-crypto/signature"

	"github.com/consensys/gnark/std/signature/eddsa"

	"github.com/akakou/zk-ban/commit"
	"github.com/akakou/zk-ban/snark"
	"github.com/consensys/gnark/frontend"
)

type SignCircuit struct {
	Credential     eddsa.Signature   `gnark:"cert"`
	UserSecretKey  frontend.Variable `gnark:"sk"`
	GroupPublicKey eddsa.PublicKey   `gnark:",public"`
	Basename       frontend.Variable `gnark:",public"`
	Message        frontend.Variable `gnark:",public"`
	Commit1        frontend.Variable `gnark:",public"`
	Commit2        frontend.Variable `gnark:",public"`
	Period         frontend.Variable `gnark:",public"`
}

func (circuit *SignCircuit) Define(api frontend.API) error {
	err := authCert(
		api,
		circuit.Period,
		circuit.UserSecretKey,
		circuit.Credential,
		circuit.GroupPublicKey)

	if err != nil {
		return err
	}

	commit1, err := snark.CircuitHash(api, circuit.Message, circuit.UserSecretKey)
	if err != nil {
		return err
	}

	commit2, err := snark.CircuitHash(api, circuit.Period, circuit.Basename, circuit.UserSecretKey)
	if err != nil {
		return err
	}

	api.AssertIsEqual(circuit.Commit1, commit1)
	api.AssertIsEqual(circuit.Commit2, commit2)

	return nil
}

func NewSignWitness(m, bsn *big.Int, commit *commit.SignCommit, signer *commit.Signer, gpk signature.PublicKey) *SignCircuit {
	assign := &SignCircuit{
		UserSecretKey: signer.UserSecretKey.Number,
		Basename:      bsn,
		Message:       m,
		Commit1:       commit.Commit1,
		Commit2:       commit.Commit2,
		Period:        signer.Period,
	}

	assign.GroupPublicKey.Assign(snark.TwistededwardsCurve, gpk.Bytes())
	assign.Credential.Assign(snark.TwistededwardsCurve, signer.Credential.Signature)

	return assign
}
