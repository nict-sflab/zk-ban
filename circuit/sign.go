package circuit

import (
	"github.com/consensys/gnark/std/signature/eddsa"

	"github.com/akakou/zk-ban/primitives"
	"github.com/akakou/zk-ban/snark"
	zkbanw "github.com/akakou/zk-ban/witness"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/frontend"
)

const MAX = 10

type SignCircuit struct {
	UserSecretKey  frontend.Variable `gnark:",secret"`
	Credential     eddsa.Signature   `gnark:",secret"`
	Random         frontend.Variable `gnark:",secret"`
	Counter        frontend.Variable `gnark:",secret"`
	Nym            frontend.Variable `gnark:",public"`
	Signature      frontend.Variable `gnark:",public"`
	Message        frontend.Variable `gnark:",public"`
	Period         frontend.Variable `gnark:",public"`
	GroupPublicKey eddsa.PublicKey   `gnark:",public"`
}

func SessionTag(api frontend.API, period, counter frontend.Variable) frontend.Variable {
	periodBits := api.ToBinary(period, 64)
	counterBits := api.ToBinary(counter, 64)

	return api.FromBinary(append(counterBits, periodBits...)...)
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

	signature, err := snark.CircuitHash(api, circuit.Message, circuit.Random, circuit.UserSecretKey)
	if err != nil {
		return err
	}

	sessionTag := SessionTag(api, circuit.Period, circuit.Counter)

	nym, err := snark.CircuitHash(api, sessionTag, circuit.UserSecretKey)
	if err != nil {
		return err
	}

	api.AssertIsLessOrEqual(circuit.Counter, MAX)

	api.AssertIsEqual(circuit.Signature, signature)
	api.AssertIsEqual(circuit.Nym, nym)

	return nil
}

func NewSignWitness(m, r *primitives.BigInt, counter int64, commit *zkbanw.SignCommit, signer *zkbanw.Signer, gpk *zkbanw.GroupPublicKey) (witness.Witness, error) {
	assign := &SignCircuit{
		UserSecretKey: signer.UserSecretKey.Number.Int,
		Period:        signer.Period,
		Counter:       counter,
		Message:       m.Int,
		Random:        r.Int,
		Signature:     commit.Sigma.Int,
		Nym:           commit.Nym.Int,
	}

	assign.GroupPublicKey.Assign(snark.TwistededwardsCurve, gpk.Bytes())
	assign.Credential.Assign(snark.TwistededwardsCurve, signer.Credential.Signature)

	return frontend.NewWitness(assign, snark.EcCurve.ScalarField())
}

func NewPublicSignWitness(m *primitives.BigInt, counter int64, commit *zkbanw.SignCommit, period int64, gpk *zkbanw.GroupPublicKey) (witness.Witness, error) {
	wit, err := NewSignWitness(m, primitives.NewBigInt(0), counter, commit, &zkbanw.Signer{
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
