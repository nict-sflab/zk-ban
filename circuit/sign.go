package circuit

import (
	"math/big"

	"github.com/consensys/gnark-crypto/signature"

	"github.com/consensys/gnark/std/signature/eddsa"

	"github.com/akakou/zk-ban/snark"
	zkbanw "github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/frontend"
)

type SignCircuit struct {
	UserPublicKey frontend.Variable `gnark:"pk"`
	Certificate   eddsa.Signature   `gnark:"cert"`
	Nonce         frontend.Variable `gnark:"nonce"`

	UserSecretKey  frontend.Variable `gnark:"sk"`
	GroupPublicKey eddsa.PublicKey   `gnark:",public"`
	Basename       frontend.Variable `gnark:",public"`
	Message        frontend.Variable `gnark:",public"`
	Commit1        frontend.Variable `gnark:",public"`
	Commit2        frontend.Variable `gnark:",public"`
	Commit3        frontend.Variable `gnark:",public"`
	Period         frontend.Variable `gnark:",public"`
}

func (circuit *SignCircuit) Define(api frontend.API) error {
	err := auth(
		api,
		circuit.Period,
		circuit.UserSecretKey,
		circuit.UserPublicKey,
		circuit.Certificate,
		circuit.GroupPublicKey)

	if err != nil {
		return err
	}

	commit1, err := hash(api, circuit.Basename, circuit.UserSecretKey)
	if err != nil {
		return err
	}

	commit2, err := hash(api, circuit.Message, circuit.Nonce)
	if err != nil {
		return err
	}

	commit3, err := hash(api, commit2, circuit.UserSecretKey)
	if err != nil {
		return err
	}

	api.AssertIsEqual(circuit.Commit1, commit1)
	api.AssertIsEqual(circuit.Commit2, commit2)
	api.AssertIsEqual(circuit.Commit3, commit3)

	return nil
}

func NewSignWitness(m, bsn, nonce *big.Int, commit *zkbanw.SignCommit, signer *zkbanw.Signer, gpk signature.PublicKey) *SignCircuit {
	assign := &SignCircuit{
		Nonce:         nonce,
		UserSecretKey: signer.UserSecretKey.Number,
		Message:       m,
		Commit1:       commit.Commit1,
		Commit2:       commit.Commit2,
		Commit3:       commit.Commit3,
		Basename:      bsn,
		Period:        signer.Period,
		UserPublicKey: signer.UserPublicKey.Number,
	}

	assign.GroupPublicKey.Assign(snark.TwistededwardsCurve, gpk.Bytes())
	assign.Certificate.Assign(snark.TwistededwardsCurve, signer.Certificate.Signature)

	return assign
}
