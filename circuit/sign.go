package circuit

import (
	"github.com/consensys/gnark/std/signature/eddsa"

	"github.com/akakou/zk-ban/primitives"
	"github.com/akakou/zk-ban/snark"
	zkbanw "github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/frontend"
)

type SignCircuit struct {
	UserSecretKey  frontend.Variable `gnark:",secret"`
	Credential     eddsa.Signature   `gnark:",secret"`
	SessionTag     frontend.Variable `gnark:",public"`
	Nym            frontend.Variable `gnark:",public"`
	Signature      frontend.Variable `gnark:",public"`
	Message        frontend.Variable `gnark:",public"`
	Period         frontend.Variable `gnark:",public"`
	GroupPublicKey eddsa.PublicKey   `gnark:",public"`
}

func (circuit *SignCircuit) Define(api frontend.API) error {
	err := authCredential(
		api,
		circuit.UserSecretKey,
		circuit.Credential,
		circuit.Period,
		circuit.GroupPublicKey)

	if err != nil {
		return err
	}

	signature, err := snark.CircuitHash(api, circuit.Message, circuit.UserSecretKey)
	if err != nil {
		return err
	}

	nym, err := snark.CircuitHash(api, circuit.SessionTag, circuit.UserSecretKey)
	if err != nil {
		return err
	}

	api.AssertIsEqual(circuit.Signature, signature)
	api.AssertIsEqual(circuit.Nym, nym)

	return nil
}

func NewSignWitness(m, sessionTag *primitives.BigInt, commit *zkbanw.SignCommit, signer *zkbanw.Signer, gpk *zkbanw.GroupPublicKey) (witness.Witness, error) {
	assign := &SignCircuit{
		UserSecretKey: signer.UserSecretKey.Number.Int,
		Period:        signer.Period,
		SessionTag:    sessionTag.Int,
		Message:       m.Int,
		Signature:     commit.Sigma.Int,
		Nym:           commit.Nym.Int,
	}

	assign.GroupPublicKey.Assign(snark.TwistededwardsCurve, gpk.Bytes())
	assign.Credential.Assign(snark.TwistededwardsCurve, signer.Credential.Signature)

	return frontend.NewWitness(assign, snark.EcCurve.ScalarField())
}

func NewPublicSignWitness(m, sessionTag *primitives.BigInt, commit *zkbanw.SignCommit, period int64, gpk *zkbanw.GroupPublicKey) (witness.Witness, error) {
	wit, err := NewSignWitness(m, sessionTag, commit, &zkbanw.Signer{
		UserSecretKey: &zkbanw.UserSecretKey{
			primitives.NewBigInt(0),
		},
		Credential: &zkbanw.Credential{
			Signature: make([]byte, 32),
		},
		Period: period,
	}, gpk)
	if err != nil {
		return nil, err
	}

	return wit.Public()
}
