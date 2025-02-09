package circuit

import (
	"math/big"

	"github.com/consensys/gnark/std/signature/eddsa"

	"github.com/akakou/zk-ban/commit"
	"github.com/akakou/zk-ban/snark"
	"github.com/consensys/gnark/frontend"
)

type SignCircuit struct {
	UserSecretKey      frontend.Variable `gnark:",secret"`
	CredentialAuthInfo CredentialAuthInfo
	GroupPublicKey     eddsa.PublicKey   `gnark:",public"`
	Basename           frontend.Variable `gnark:",public"`
	Message            frontend.Variable `gnark:",public"`
	Signature          frontend.Variable `gnark:",public"`
	Nym                frontend.Variable `gnark:",public"`
}

func (circuit *SignCircuit) Define(api frontend.API) error {
	err := authCredential(
		api,
		circuit.CredentialAuthInfo,
		circuit.UserSecretKey,
		circuit.GroupPublicKey)

	if err != nil {
		return err
	}

	signature, err := snark.CircuitHash(api, circuit.Message, circuit.UserSecretKey)
	if err != nil {
		return err
	}

	nym, err := snark.CircuitHash(api, circuit.CredentialAuthInfo.Period, circuit.Basename, circuit.UserSecretKey)
	if err != nil {
		return err
	}

	api.AssertIsEqual(circuit.Signature, signature)
	api.AssertIsEqual(circuit.Nym, nym)

	return nil
}

func NewSignWitness(m, bsn *big.Int, commit *commit.SignCommit, signer *commit.Signer) *SignCircuit {
	assign := &SignCircuit{
		UserSecretKey: signer.UserSecretKey.Number,
		CredentialAuthInfo: CredentialAuthInfo{
			Period: signer.Period,
		},
		Basename:  bsn,
		Message:   m,
		Signature: commit.Commit1,
		Nym:       commit.Commit2,
	}

	assign.GroupPublicKey.Assign(snark.TwistededwardsCurve, signer.GroupPublicKey.Bytes())
	assign.CredentialAuthInfo.Credential.Assign(snark.TwistededwardsCurve, signer.Credential.Signature)

	return assign
}
